type Props = { className?: string };

export default function Logo({ className = "h-7 w-7" }: Props) {
  return (
    <svg
      viewBox="0 0 32 32"
      fill="none"
      className={className}
      aria-hidden
    >
      <rect width="32" height="32" rx="7" fill="#0c0c0e" />
      <rect x="0.5" y="0.5" width="31" height="31" rx="6.5" stroke="#27272a" />
      <path
        d="M9 11.5C9 10.1193 10.1193 9 11.5 9H20.5C21.8807 9 23 10.1193 23 11.5V20.5C23 21.8807 21.8807 23 20.5 23H11.5C10.1193 23 9 21.8807 9 20.5V11.5Z"
        stroke="#34d399"
        strokeWidth="1.6"
      />
      <circle cx="16" cy="16" r="3" fill="#34d399" />
      <path
        d="M16 4V7M16 25V28M4 16H7M25 16H28"
        stroke="#34d399"
        strokeWidth="1.6"
        strokeLinecap="round"
      />
    </svg>
  );
}
