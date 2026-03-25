"use client";

export default function TeamDetailError() {
  return (
    <div className="py-10">
      <h1 className="text-2xl font-semibold">Could not load team details</h1>
      <p className="mt-2 text-sm text-neutral-600">The public team page is temporarily unavailable.</p>
    </div>
  );
}
