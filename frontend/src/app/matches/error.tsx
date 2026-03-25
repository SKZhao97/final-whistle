"use client";

export default function MatchesError() {
  return (
    <div className="py-10">
      <h1 className="text-2xl font-semibold">Could not load matches</h1>
      <p className="mt-2 text-sm text-neutral-600">
        The public match list is temporarily unavailable.
      </p>
    </div>
  );
}
