export default function MatchNotFound() {
  return (
    <div className="py-10">
      <h1 className="text-2xl font-semibold">Match not found</h1>
      <p className="mt-2 text-sm text-neutral-600">
        The requested match does not exist or is not available.
      </p>
    </div>
  );
}
