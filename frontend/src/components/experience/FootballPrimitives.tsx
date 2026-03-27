import type { ReactNode } from "react";

import type { TeamSummary } from "@/types/api";

type SectionShellProps = {
  eyebrow?: string;
  title: string;
  description?: string;
  children: ReactNode;
  accent?: "field" | "paper";
  className?: string;
};

export function SectionShell({
  eyebrow,
  title,
  description,
  children,
  accent = "paper",
  className = "",
}: SectionShellProps) {
  return (
    <section
      className={`match-shell ${accent === "field" ? "match-shell--field" : ""} ${className}`.trim()}
    >
      <header className="space-y-2">
        {eyebrow ? <p className="match-eyebrow">{eyebrow}</p> : null}
        <div className="space-y-1">
          <h2 className="text-xl font-semibold tracking-tight text-[var(--fw-ink)]">{title}</h2>
          {description ? (
            <p className="max-w-3xl text-sm leading-6 text-[var(--fw-muted)]">{description}</p>
          ) : null}
        </div>
      </header>
      {children}
    </section>
  );
}

export function BrandMark() {
  return (
    <div className="brand-mark" aria-hidden="true">
      <span className="brand-mark__outer" />
      <span className="brand-mark__mid" />
      <span className="brand-mark__arc" />
      <span className="brand-mark__arc brand-mark__arc--bottom" />
      <span className="brand-mark__spot" />
    </div>
  );
}

export function Wordmark() {
  return (
    <div className="space-y-0.5 leading-none">
      <p className="text-[0.65rem] font-semibold uppercase tracking-[0.28em] text-[var(--fw-field-700)]">
        Match Archive
      </p>
      <p className="font-semibold tracking-[-0.04em] text-[var(--fw-ink)] sm:text-[1.1rem]">
        Final <span className="text-[var(--fw-field-700)]">Whistle</span>
      </p>
    </div>
  );
}

export function TeamCrest({
  team,
  size = "lg",
}: {
  team: TeamSummary;
  size?: "sm" | "md" | "lg";
}) {
  const sizeClass =
    size === "sm"
      ? "h-12 w-12"
      : size === "md"
        ? "h-16 w-16"
        : "h-20 w-20";

  const palette = teamPalette[team.slug] ?? teamPalette.default;
  const shouldUseRemoteLogo = team.logoUrl && !team.logoUrl.includes("example.com");

  return (
    <div
      className={`crest-shell ${sizeClass}`}
      style={{
        background: shouldUseRemoteLogo
          ? undefined
          : `linear-gradient(180deg, ${palette[0]}, ${palette[1]})`,
      }}
    >
      {shouldUseRemoteLogo ? (
        // eslint-disable-next-line @next/next/no-img-element
        <img src={team.logoUrl} alt={team.name} className="h-full w-full object-contain" />
      ) : (
        <span className="text-sm font-semibold uppercase tracking-[0.18em] text-white">
          {team.shortName ?? team.name.slice(0, 3)}
        </span>
      )}
    </div>
  );
}

export function LeagueMark({ label }: { label: string }) {
  return (
    <div className="inline-flex items-center gap-2 rounded-full border border-[var(--fw-line)] bg-[var(--fw-surface)]/82 px-3 py-1.5 text-xs font-medium text-[var(--fw-muted)] shadow-[0_12px_30px_rgba(16,31,24,0.05)]">
      <span className="league-mark-icon">PL</span>
      <span>{label}</span>
    </div>
  );
}

export function ArchivePill({ children }: { children: ReactNode }) {
  return (
    <span className="inline-flex items-center rounded-full border border-[var(--fw-line)] bg-[var(--fw-surface)] px-3 py-1 text-xs font-medium text-[var(--fw-ink-soft)]">
      {children}
    </span>
  );
}

export function ArchiveStat({
  label,
  value,
  detail,
}: {
  label: string;
  value: string;
  detail?: string;
}) {
  return (
    <div className="rounded-[1.5rem] border border-[var(--fw-line)] bg-[var(--fw-surface)] p-5 shadow-[0_18px_40px_rgba(16,31,24,0.08)]">
      <p className="text-xs uppercase tracking-[0.18em] text-[var(--fw-muted)]">{label}</p>
      <p className="mt-3 text-2xl font-semibold tracking-tight text-[var(--fw-ink)]">{value}</p>
      {detail ? <p className="mt-2 text-sm text-[var(--fw-muted)]">{detail}</p> : null}
    </div>
  );
}

const teamPalette: Record<string, [string, string]> = {
  "manchester-city": ["#69a8ff", "#12436b"],
  liverpool: ["#d44c62", "#7b1726"],
  arsenal: ["#d74a46", "#8f1d1f"],
  chelsea: ["#4f7fe2", "#153786"],
  "manchester-united": ["#dd6150", "#7e2118"],
  "tottenham-hotspur": ["#35527b", "#12283f"],
  default: ["#406a54", "#1d3528"],
};
