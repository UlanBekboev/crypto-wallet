"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");
  const [messageColor, setMessageColor] = useState("text-green-600");
  const [loading, setLoading] = useState(false);
  const [blockedTime, setBlockedTime] = useState<number | null>(null);

  useEffect(() => {
    if (!blockedTime) return;

    const interval = setInterval(() => {
      setBlockedTime((prev) => (prev && prev > 1 ? prev - 1 : null));
    }, 1000);

    return () => clearInterval(interval);
  }, [blockedTime]);


  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setMessage("");

    try {
      const res = await fetch("http://localhost:8080/api/auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ email, password }),
      });

      const data = await res.json();

      if (res.status === 429 && data.blocked_secs) {
        setBlockedTime(data.blocked_secs);
        setMessage("Слишком много попыток. Подождите.");
        setMessageColor("text-red-600");
        return;
      }


      if (res.ok) {
        setMessage("Успешный вход. Перенаправление...");
        setMessageColor("text-green-600");
        setTimeout(() => router.push("/profile"), 1500);
      } else {
        setMessage(data.error || "Ошибка входа");
        setMessageColor("text-red-600");
      }
    } catch (error) {
      setMessage("Сервер недоступен.");
      setMessageColor("text-red-600");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="max-w-md mx-auto mt-20 space-y-4 bg-white p-6 rounded shadow">
      <h2 className="text-xl font-bold">Вход</h2>
      {message && <p className={`font-medium ${messageColor}`}>{message}</p>}
      <input
        type="email"
        placeholder="Email"
        className="w-full border p-2 rounded"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        required
      />
      <input
        type="password"
        placeholder="Пароль"
        className="w-full border p-2 rounded"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        required
      />
      <button
        type="submit"
        disabled={loading}
        className={`w-full bg-blue-600 text-white py-2 px-4 rounded ${loading ? "opacity-50 cursor-not-allowed" : ""}`}
      >
        {loading ? "Вход..." : "Войти"}
      </button>

      <p className="text-sm text-center mt-2">
        <a href="/forgot-password" className="text-blue-600 hover:underline">
          Забыли пароль?
        </a>
      </p>
      <p className="text-sm text-center mt-2">
        <a href="/register" className="text-blue-600 hover:underline">
          Зарегистрироваться
        </a>
      </p>
      {blockedTime && (
        <p className="text-red-500 font-medium">
          Повторите попытку через {blockedTime} сек.
        </p>
      )}

    </form>
  );
}
