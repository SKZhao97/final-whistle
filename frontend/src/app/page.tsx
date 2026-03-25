export default function Home() {
  return (
    <div className="flex flex-col items-center justify-center py-12">
      <div className="max-w-2xl mx-auto text-center space-y-8">
        <h1 className="text-4xl font-bold tracking-tight text-gray-900 dark:text-gray-100 sm:text-5xl">
          Welcome to Final Whistle
        </h1>
        <p className="text-lg text-gray-600 dark:text-gray-400">
          A post-match recording product for football viewers.
          Record your watched matches, rate teams and players,
          and build your personal football memory archive.
        </p>
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <a
            href="/matches"
            className="inline-flex items-center justify-center rounded-md bg-primary px-6 py-3 text-sm font-medium text-primary-foreground shadow transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          >
            Browse Matches
          </a>
          <a
            href="/me"
            className="inline-flex items-center justify-center rounded-md border border-input bg-background px-6 py-3 text-sm font-medium shadow-sm transition-colors hover:bg-accent hover:text-accent-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          >
            My Profile
          </a>
        </div>
        <div className="pt-8 grid grid-cols-1 sm:grid-cols-2 gap-6 text-left">
          <div className="space-y-3">
            <h3 className="text-lg font-semibold">Record Your Matches</h3>
            <p className="text-sm text-gray-500">
              Mark matches you&apos;ve watched, note how you watched them, and when.
            </p>
          </div>
          <div className="space-y-3">
            <h3 className="text-lg font-semibold">Rate & Review</h3>
            <p className="text-sm text-gray-500">
              Score the match, both teams, and standout players. Add short reviews and tags.
            </p>
          </div>
          <div className="space-y-3">
            <h3 className="text-lg font-semibold">Build Your Archive</h3>
            <p className="text-sm text-gray-500">
              See your watch history, average ratings, and favorite teams/players.
            </p>
          </div>
          <div className="space-y-3">
            <h3 className="text-lg font-semibold">Community Perspective</h3>
            <p className="text-sm text-gray-500">
              Check aggregated ratings, player rankings, and recent reviews for each match.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
