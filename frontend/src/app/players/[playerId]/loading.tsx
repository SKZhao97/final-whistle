export default function PlayerDetailLoading() {
  return (
    <div className="py-8">
      <div className="h-10 w-60 animate-pulse rounded bg-neutral-200" />
      <div className="mt-8 rounded-xl border p-5">
        <div className="h-5 w-40 animate-pulse rounded bg-neutral-200" />
        <div className="mt-4 h-32 animate-pulse rounded bg-neutral-100" />
      </div>
    </div>
  );
}
