"use client";

import { useState } from "react";

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState("");
  const [message, setMessage] = useState("");
  const [messageColor, setMessageColor] = useState("text-green-600");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setMessage("");

    try {
      const res = await fetch("http://localhost:8080/api/auth/forgot-password", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ email }),
      });

      const data = await res.json();

      if (res.ok) {
        setMessage("Письмо отправлено на почту");
        setMessageColor("text-green-600");
        setEmail("");
      } else {
        setMessage(data.error || "Ошибка при отправке");
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
      <h2 className="text-xl font-bold">Сброс пароля</h2>
      {message && <p className={`font-medium ${messageColor}`}>{message}</p>}
      <input
        type="email"
        placeholder="Email"
        className="w-full border p-2 rounded"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        required
      />
      <button
        type="submit"
        disabled={loading}
        className={`w-full bg-blue-600 text-white py-2 px-4 rounded ${
          loading ? "opacity-50 cursor-not-allowed" : ""
        }`}
      >
        {loading ? "Отправка..." : "Отправить ссылку для сброса"}
      </button>
    </form>
  );
}
