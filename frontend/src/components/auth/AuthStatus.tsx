"use client";

import Link from "next/link";

import { useAuth } from "@/components/auth/AuthProvider";

export default function AuthStatus() {
  const { status, user, logout } = useAuth();

  if (status === "loading") {
    return <span className="text-sm text-neutral-500">Checking session...</span>;
  }

  if (status === "authenticated" && user) {
    return (
      <div className="flex items-center gap-3">
        <span className="text-sm text-neutral-600">{user.name}</span>
        <Link
          href="/me"
          className="inline-flex items-center justify-center rounded-md border px-3 py-2 text-sm font-medium transition-colors hover:bg-neutral-50"
        >
          My Profile
        </Link>
        <button
          type="button"
          onClick={() => void logout()}
          className="inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
        >
          Log Out
        </button>
      </div>
    );
  }

  return (
    <Link
      href="/login"
      className="inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
    >
      Dev Login
    </Link>
  );
}
