"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

import { useAuth } from "@/components/auth/AuthProvider";
import { ApiError } from "@/lib/api/client";

export default function LoginPage() {
  const router = useRouter();
  const { login, status } = useAuth();
  const [email, setEmail] = useState("demo@final-whistle.test");
  const [name, setName] = useState("Demo User");
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitting(true);
    setError(null);

    try {
      await login({ email, name });
      router.push("/me");
      router.refresh();
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError("Login failed");
      }
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="mx-auto max-w-md py-12">
      <p className="text-xs uppercase tracking-[0.2em] text-neutral-500">Session Auth</p>
      <h1 className="mt-2 text-3xl font-bold tracking-tight">Development Login</h1>
      <p className="mt-3 text-sm text-neutral-600">
        Use a seeded dev user or submit a new email while development auto-create is enabled.
      </p>

      <form onSubmit={handleSubmit} className="mt-8 space-y-4 rounded-2xl border p-6">
        <label className="block text-sm">
          <span className="mb-2 block font-medium">Email</span>
          <input
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            className="w-full rounded-md border px-3 py-2"
            type="email"
            required
          />
        </label>

        <label className="block text-sm">
          <span className="mb-2 block font-medium">Display Name</span>
          <input
            value={name}
            onChange={(event) => setName(event.target.value)}
            className="w-full rounded-md border px-3 py-2"
            type="text"
            required
          />
        </label>

        {error ? <p className="text-sm text-red-600">{error}</p> : null}
        {status === "authenticated" ? (
          <p className="text-sm text-emerald-700">You are already signed in.</p>
        ) : null}

        <button
          type="submit"
          disabled={submitting}
          className="inline-flex w-full items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 disabled:opacity-50"
        >
          {submitting ? "Signing In..." : "Sign In"}
        </button>
      </form>
    </div>
  );
}
