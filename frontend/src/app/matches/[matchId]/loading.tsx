export default function MatchDetailLoading() {
  return (
    <div className="py-8">
      <div className="mb-4 h-6 w-40 animate-pulse rounded bg-neutral-200" />
      <div className="mb-8 h-10 w-3/4 animate-pulse rounded bg-neutral-100" />
      <div className="grid gap-4 lg:grid-cols-3">
        {Array.from({ length: 3 }).map((_, index) => (
          <div key={index} className="rounded-xl border p-5">
            <div className="h-5 w-32 animate-pulse rounded bg-neutral-200" />
            <div className="mt-4 h-20 animate-pulse rounded bg-neutral-100" />
          </div>
        ))}
      </div>
    </div>
  );
}
