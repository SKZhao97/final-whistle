"use client";

/**
 * 认证提供者组件和钩子。
 * 提供用户认证状态管理和操作。
 */

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";

import { ApiError, authApi } from "@/lib/api/client";
import type { LoginRequest, UserSummary } from "@/types/api";

/** 认证状态类型。 */
type AuthStatus = "loading" | "authenticated" | "unauthenticated";

/** 认证上下文值类型。 */
type AuthContextValue = {
  /** 当前认证状态。 */
  status: AuthStatus;
  /** 当前用户信息，未认证时为 null。 */
  user: UserSummary | null;
  /** 登录函数。 */
  login: (input: LoginRequest) => Promise<void>;
  /** 登出函数。 */
  logout: () => Promise<void>;
  /** 刷新认证状态函数。 */
  refresh: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

/**
 * 认证提供者组件。
 *
 * 管理用户认证状态，提供登录、登出、刷新状态等功能。
 * 自动恢复会话，监听认证状态变化。
 *
 * @param children - 子组件
 */
export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [status, setStatus] = useState<AuthStatus>("loading");
  const [user, setUser] = useState<UserSummary | null>(null);

  const refresh = useCallback(async () => {
    try {
      const result = await authApi.me();
      setUser(result.user);
      setStatus("authenticated");
    } catch (error) {
      if (error instanceof ApiError && error.code === "UNAUTHORIZED") {
        setUser(null);
        setStatus("unauthenticated");
        return;
      }
      throw error;
    }
  }, []);

  useEffect(() => {
    let cancelled = false;

    async function restoreSession() {
      try {
        const result = await authApi.me();
        if (cancelled) {
          return;
        }
        setUser(result.user);
        setStatus("authenticated");
      } catch (error) {
        if (cancelled) {
          return;
        }
        if (error instanceof ApiError && error.code === "UNAUTHORIZED") {
          setUser(null);
          setStatus("unauthenticated");
          return;
        }
        setUser(null);
        setStatus("unauthenticated");
      }
    }

    void restoreSession();
    return () => {
      cancelled = true;
    };
  }, []);

  const login = useCallback(async (input: LoginRequest) => {
    const result = await authApi.login(input);
    setUser(result.user);
    setStatus("authenticated");
  }, []);

  const logout = useCallback(async () => {
    await authApi.logout();
    setUser(null);
    setStatus("unauthenticated");
  }, []);

  const value = useMemo(
    () => ({ status, user, login, logout, refresh }),
    [status, user, login, logout, refresh],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

/**
 * 使用认证上下文钩子。
 *
 * 在组件中获取认证状态和操作函数。
 * 必须在 AuthProvider 内部使用。
 *
 * @returns 认证上下文值
 * @throws 如果在 AuthProvider 外部使用则抛出错误
 */
export function useAuth() {
  const value = useContext(AuthContext);
  if (!value) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return value;
}
