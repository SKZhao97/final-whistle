import Link from "next/link";

export default function TeamsPage() {
  return (
    <div className="py-8">
      <h1 className="text-3xl font-bold tracking-tight">Teams</h1>
      <p className="mt-3 text-sm text-neutral-600">
        Team detail pages are currently entered from match detail pages. Browse matches first to open a team.
      </p>
      <Link
        href="/matches"
        className="mt-6 inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
      >
        Browse Matches
      </Link>
    </div>
  );
}
