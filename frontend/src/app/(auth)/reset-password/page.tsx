"use client";

import { useSearchParams, useRouter } from "next/navigation";
import { useState } from "react";

export default function ResetPasswordPage() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const token = searchParams.get("token");

  const [newPassword, setNewPassword] = useState("");
  const [message, setMessage] = useState("");
  const [messageColor, setMessageColor] = useState("text-green-600");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token) {
      setMessage("Отсутствует токен");
      setMessageColor("text-red-600");
      return;
    }

    setLoading(true);
    setMessage("");

    try {
      const res = await fetch("http://localhost:8080/api/auth/reset-password", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ token, newPassword }),
      });

      const data = await res.json();

      if (res.ok) {
        setMessage("Пароль успешно сброшен");
        setMessageColor("text-green-600");
        setNewPassword("");

        setTimeout(() => {
          router.push("/login");
        }, 2000);
      } else {
        setMessage(data.error || "Ошибка сброса");
        setMessageColor("text-red-600");
      }
    } catch (error) {
      setMessage("Ошибка соединения с сервером");
      setMessageColor("text-red-600");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="max-w-md mx-auto mt-20 bg-white p-6 rounded shadow space-y-4">
      <h2 className="text-xl font-bold">Новый пароль</h2>
      {message && <p className={`font-medium ${messageColor}`}>{message}</p>}
      <input
        type="password"
        placeholder="Введите новый пароль"
        className="w-full border p-2 rounded"
        value={newPassword}
        onChange={(e) => setNewPassword(e.target.value)}
        required
      />
      <button
        type="submit"
        disabled={loading}
        className={`w-full bg-blue-600 text-white py-2 px-4 rounded ${
          loading ? "opacity-50 cursor-not-allowed" : ""
        }`}
      >
        {loading ? "Отправка..." : "Сбросить пароль"}
      </button>
    </form>
  );
}
