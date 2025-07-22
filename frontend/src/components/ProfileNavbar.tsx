"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const links = [
  { href: "/wallet", label: "💰 Кошелёк" },
  { href: "/transfer", label: "🔁 Перевод" },
  { href: "/history", label: "📜 История" },
];

export function ProfileNavbar() {
  const pathname = usePathname();

  return (
    <nav className="flex justify-center gap-4 py-4">
      {links.map(({ href, label }) => {
        const isActive = pathname === href;
        return (
          <Link
            key={href}
            href={href}
            className={`px-4 py-2 rounded transition ${
              isActive
                ? "bg-blue-600 text-white"
                : "bg-gray-100 text-gray-800 hover:bg-blue-100"
            }`}
          >
            {label}
          </Link>
        );
      })}
    </nav>
  );
}
