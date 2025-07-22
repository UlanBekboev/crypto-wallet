"use client";

import { useState } from "react";
import { fetchWithRefresh } from "@/lib/fetchWithRefresh";
import Link from "next/link";

export default function TransferPage() {
  const [recipient, setRecipient] = useState("");
  const [amount, setAmount] = useState("");
  const [message, setMessage] = useState("");

  const handleTransfer = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      const res = await fetchWithRefresh("http://localhost:8080/api/auth/transfer", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ to: recipient, amount: parseFloat(amount) }),
      });

      if (res.ok) {
        setMessage("✅ Перевод выполнен!");
        setRecipient("");
        setAmount("");
      } else {
        const data = await res.json();
        setMessage(`❌ Ошибка: ${data.message || "неизвестная ошибка"}`);
      }
    } catch (err) {
      setMessage("❌ Ошибка отправки запроса");
    }
  };

  return (
    <div className="max-w-md mx-auto mt-10 p-6 bg-white rounded shadow">
      <h1 className="text-2xl font-bold mb-4 text-center">🔁 Перевод</h1>

      <form onSubmit={handleTransfer} className="space-y-4">
        <input
          type="text"
          placeholder="ID получателя"
          value={recipient}
          onChange={(e) => setRecipient(e.target.value)}
          className="w-full border p-2 rounded"
          required
        />
        <input
          type="number"
          placeholder="Сумма"
          value={amount}
          onChange={(e) => setAmount(e.target.value)}
          className="w-full border p-2 rounded"
          required
        />
        <button
          type="submit"
          className="w-full bg-green-600 text-white px-4 py-2 rounded"
        >
          Отправить
        </button>
      </form>

      {message && <p className="mt-4 text-center">{message}</p>}

      <Link
        href="/profile"
        className="mt-6 inline-block bg-blue-600 text-white px-4 py-2 rounded"
      >
        ← Назад в профиль
      </Link>
    </div>
  );
}
