"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

export default function RegisterPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [name, setName] = useState("");
  const [message, setMessage] = useState("");
  const [messageColor, setMessageColor] = useState("text-green-600");
  const [loading, setLoading] = useState(false);

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setMessage("");

    try {
      const res = await fetch("http://localhost:8080/api/auth/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ email, password, name }),
      });

      const data = await res.json();

      if (res.ok) {
        setMessage("Регистрация успешна. Перенаправление...");
        setMessageColor("text-green-600");
        setTimeout(() => router.push("/login"), 1500);
      } else {
        setMessage(data.error || "Ошибка регистрации");
        setMessageColor("text-red-600");
      }
    } catch (err) {
      setMessage("Ошибка соединения с сервером");
      setMessageColor("text-red-600");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleRegister} className="max-w-md mx-auto mt-20 space-y-4 bg-white p-6 rounded shadow">
      <h2 className="text-xl font-bold">Регистрация</h2>
      {message && <p className={`font-medium ${messageColor}`}>{message}</p>}
      <input
        type="text"
        placeholder="Имя"
        className="w-full border p-2 rounded"
        value={name}
        onChange={(e) => setName(e.target.value)}
        required
      />
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
        className={`w-full bg-blue-600 text-white py-2 px-4 rounded ${
          loading ? "opacity-50 cursor-not-allowed" : ""
        }`}
      >
        {loading ? "Регистрация..." : "Зарегистрироваться"}
      </button>
      <p className="text-sm text-center mt-2">
        <a href="/login" className="text-blue-600 hover:underline">
          Уже есть аккаунт? Войти
        </a>
      </p>
    </form>
  );
}
