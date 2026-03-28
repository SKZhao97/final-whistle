"use client";

import { useRouter } from "next/navigation";

type SeasonOption = {
  id: string;
  label: string;
};

type SeasonSelectorProps = {
  label: string;
  value: string;
  options: SeasonOption[];
};

export function SeasonSelector({ label, value, options }: SeasonSelectorProps) {
  const router = useRouter();

  return (
    <div className="flex items-center justify-end gap-2">
      <label
        htmlFor="season"
        className="text-[11px] font-semibold uppercase tracking-[0.2em] text-[var(--fw-muted)]"
      >
        {label}
      </label>
      <select
        id="season"
        name="season"
        value={value}
        onChange={(event) => {
          const season = event.target.value;
          if (options[0]?.id === season) {
            router.push("/matches");
            return;
          }
          router.push(`/matches?season=${encodeURIComponent(season)}`);
        }}
        className="min-w-[10rem] rounded-full border border-[var(--fw-line)] bg-[var(--fw-panel)] px-4 py-2 text-sm font-medium text-[var(--fw-ink)] outline-none transition focus:border-[var(--fw-field-600)]"
      >
        {options.map((season) => (
          <option key={season.id} value={season.id}>
            {season.label}
          </option>
        ))}
      </select>
    </div>
  );
}
