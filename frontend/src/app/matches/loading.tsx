export default function MatchesLoading() {
  return (
    <div className="py-8">
      <div className="mb-6 h-8 w-48 animate-pulse rounded bg-neutral-200" />
      <div className="grid gap-4">
        {Array.from({ length: 4 }).map((_, index) => (
          <div key={index} className="rounded-xl border p-5">
            <div className="mb-3 h-5 w-40 animate-pulse rounded bg-neutral-200" />
            <div className="h-4 w-64 animate-pulse rounded bg-neutral-100" />
          </div>
        ))}
      </div>
    </div>
  );
}
