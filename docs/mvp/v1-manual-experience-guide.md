# Final Whistle V1 Manual Experience Guide

This document explains how to run the current V1 locally and experience the full product loop:

`login -> browse matches -> view match detail -> create/edit check-in -> view profile`

## 1. Prerequisites

Make sure your machine has:

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- A local PostgreSQL database named `final_whistle`

If the database does not exist yet:

```bash
createdb final_whistle
```

## 2. Start the Backend

Open a terminal:

```bash
cd /Users/sz/Code/final-whistle/backend

export DATABASE_URL="postgres://sz@localhost:5432/final_whistle?sslmode=disable"
export PORT=8080
export ENV=development

go run ./cmd/migrate/main.go
go run ./cmd/seed/main.go -scope=all
go run ./cmd/api/main.go
```

If your local PostgreSQL username is not `sz`, replace it in `DATABASE_URL`.

### Verify the backend

```bash
curl http://localhost:8080/health
```

Expected result:

```json
{"status":"ok","time":"..."}
```

## 3. Start the Frontend

Open another terminal:

```bash
cd /Users/sz/Code/final-whistle/frontend
npm install
npm run dev
```

Default frontend URL:

- `http://localhost:3000`

If your backend is not running on `http://localhost:8080`, set:

```bash
export NEXT_PUBLIC_API_URL="http://localhost:8080"
npm run dev
```

## 4. Recommended Experience Flow

Follow this exact path for the most complete V1 experience.

### Step 1: Open the home page

Visit:

- `http://localhost:3000/`

You should see the Final Whistle landing page.

### Step 2: Browse matches

Visit:

- `http://localhost:3000/matches`

Expected:

- A populated match list from seed data
- Match cards with teams, status, date, and score where available
- Working navigation into match detail pages

### Step 3: View a match detail page

Open any finished match, for example:

- `http://localhost:3000/matches/1`

Expected:

- Match base information
- Aggregate summary
- Recent reviews
- Player rating summary
- Check-in panel

If you are not logged in yet, the check-in panel should show a login prompt instead of trying to load private data.

### Step 4: Log in

Visit:

- `http://localhost:3000/login`

Recommended dev credentials:

- Email: `demo@final-whistle.test`
- Name: `Demo User`

You can also use any new development email, for example:

- `yourname@final-whistle.test`

In development mode, a new user will be created automatically.

After login, you should land on:

- `http://localhost:3000/me`

### Step 5: View the profile page

Expected on `/me`:

- Current signed-in user information
- Check-in count
- Average match rating
- Recent 30-day count
- Favorite team and most-used tag where data exists
- Check-in history list

If you use a fresh user, the expected result is:

- Most counters are `0` or empty
- No check-in history yet
- A clear empty-state prompt to browse matches

### Step 6: Create a check-in

Go back to a finished match detail page, for example:

- `http://localhost:3000/matches/1`

In the check-in panel, fill in:

- Watched Type
- Supporter Side
- Match Rating
- Home Team Rating
- Away Team Rating
- Watched At
- Tags
- Short Review
- Player Ratings

Current rules:

- Only `FINISHED` matches can be submitted
- Ratings must be integers from `1` to `10`
- Rateable players come from that match's `match_players`
- You cannot select the same player twice
- There is no longer a five-player limit

Expected after submit:

- The page refreshes
- The check-in panel changes into saved state
- You can enter edit mode

### Step 7: Edit the same check-in

Use the edit flow on the same match page and change a few values.

Expected:

- Update succeeds
- Tags and player ratings are replaced by the new values
- The saved state reflects the latest content immediately

### Step 8: Return to the profile page

Visit:

- `http://localhost:3000/me`

Expected:

- `checkInCount` has increased
- `avgMatchRating` is updated
- `recentCheckInCount` is updated
- The new match appears in history
- Favorite team / most-used tag may now appear

This step validates the most important V1 loop.

## 5. Direct API Verification

If you want to verify functionality without the frontend, use `curl`.

### Login and store cookie

```bash
curl -c /tmp/fw-cookie.txt \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@final-whistle.test","name":"Demo User"}' \
  http://localhost:8080/auth/login
```

### Current user

```bash
curl -b /tmp/fw-cookie.txt http://localhost:8080/auth/me
```

### Profile summary

```bash
curl -b /tmp/fw-cookie.txt http://localhost:8080/me/profile
```

### Check-in history

```bash
curl -b /tmp/fw-cookie.txt "http://localhost:8080/me/checkins?page=1&pageSize=10"
```

### Public match detail

```bash
curl http://localhost:8080/matches/1
```

### Current user's check-in for a match

```bash
curl -b /tmp/fw-cookie.txt http://localhost:8080/matches/1/my-checkin
```

### Create a check-in

```bash
curl -b /tmp/fw-cookie.txt \
  -H "Content-Type: application/json" \
  -d '{
    "watchedType":"FULL",
    "supporterSide":"HOME",
    "matchRating":8,
    "homeTeamRating":9,
    "awayTeamRating":7,
    "shortReview":"Great match",
    "watchedAt":"2026-03-26T10:00:00Z",
    "tags":[1,4],
    "playerRatings":[
      {"playerId":1,"rating":9},
      {"playerId":2,"rating":8}
    ]
  }' \
  http://localhost:8080/matches/1/checkin
```

## 6. What To Check During Experience

Focus on these behaviors:

- Unauthenticated users see a login prompt on match detail instead of private check-in data
- `/me` loads correctly after login
- A fresh user sees a meaningful empty state
- Creating a check-in updates both match detail and `/me`
- Editing a check-in updates saved values immediately
- Invalid input returns useful validation feedback
- Finished and non-finished matches behave differently

## 7. Known Imperfections

These are currently known rough edges:

- `/teams` and `/players` list routes are still placeholder pages
- The `/me` page can still be improved for session-expiry fallback after the user is already marked authenticated

These do not block the main V1 experience path, but they are known cleanup items.

## 8. Minimal Quick Experience Path

If you want the shortest possible manual test:

1. Start backend
2. Start frontend
3. Open `/login`
4. Log in
5. Open `/matches/1`
6. Submit a check-in
7. Open `/me`

If this path works smoothly, you have experienced the core V1 loop successfully.
