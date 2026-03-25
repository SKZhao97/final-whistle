"use client";

export default function PlayerDetailError() {
  return (
    <div className="py-10">
      <h1 className="text-2xl font-semibold">Could not load player details</h1>
      <p className="mt-2 text-sm text-neutral-600">The public player page is temporarily unavailable.</p>
    </div>
  );
}
