"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { fetchWithRefresh } from "@/lib/fetchWithRefresh";

type User = {
  id: number;
  email: string;
  name?: string;
};

export default function ProfilePage() {
  const [user, setUser] = useState<User | null>(null);
  const router = useRouter();

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const res = await fetchWithRefresh("http://localhost:8080/api/auth/me");
        if (res.ok) {
          const data = await res.json();
          setUser(data.user);
        } else {
          router.push("/login");
        }
      } catch (err) {
        console.error("Ошибка при получении профиля:", err);
        router.push("/login");
      }
    };

    fetchProfile();
  }, []);


  if (!user)
    return <div className="text-center mt-10">Загрузка профиля...</div>;

  return (
    <div className="max-w-xl mx-auto mt-20 bg-white p-6 rounded shadow space-y-4">
      <h2 className="text-2xl font-bold">Профиль пользователя</h2>
      <p>
        <strong>ID:</strong> {user.id}
      </p>
      <p>
        <strong>Email:</strong> {user.email}
      </p>
      {user.name && (
        <p>
          <strong>Имя:</strong> {user.name}
        </p>
      )}
      <button
          onClick={async () => {
            await fetch("http://localhost:8080/api/auth/logout", {
              method: "POST",
              credentials: "include",
            });
            router.push("/login");
          }}
          className="bg-red-600 text-white px-4 py-2 rounded"
        >
          Выйти
      </button>
      <Link
        href="/change-password"
        className="inline-block mt-4 bg-blue-600 text-white px-4 py-2 rounded"
      >
        Сменить пароль
      </Link>
    </div>
  );
}
