"use client";

import Link from "next/link";

import { useAuth } from "@/components/auth/AuthProvider";

export default function MePage() {
  const { status, user } = useAuth();

  if (status === "loading") {
    return (
      <div className="py-8">
        <h1 className="text-3xl font-bold">My Profile</h1>
        <p className="mt-4 text-neutral-600">Checking your session...</p>
      </div>
    );
  }

  if (status === "unauthenticated" || !user) {
    return (
      <div className="py-8">
        <h1 className="text-3xl font-bold">My Profile</h1>
        <p className="mt-4 text-neutral-600">
          You need to sign in before your profile and future check-ins are available.
        </p>
        <Link
          href="/login"
          className="mt-6 inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
        >
          Go to Dev Login
        </Link>
      </div>
    );
  }

  return (
    <div className="py-8">
      <p className="text-xs uppercase tracking-[0.2em] text-neutral-500">Authenticated</p>
      <h1 className="mt-2 text-3xl font-bold tracking-tight">My Profile</h1>
      <div className="mt-8 rounded-2xl border p-6">
        <p className="text-sm text-neutral-500">Signed in as</p>
        <p className="mt-2 text-xl font-semibold">{user.name}</p>
        <p className="mt-4 text-sm text-neutral-600">
          Profile details and check-in history will arrive in the later profile spec.
        </p>
      </div>
    </div>
  );
}
