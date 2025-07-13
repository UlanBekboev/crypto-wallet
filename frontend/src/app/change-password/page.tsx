"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";

export default function ChangePasswordPage() {
  const [oldPassword, setOldPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [message, setMessage] = useState("");
  const [messageColor, setMessageColor] = useState("text-green-600");
  const router = useRouter();

  const handleBack = () => {
    router.back(); // возвращает на предыдущую страницу
  };

  const handleChange = async (e: React.FormEvent) => {
    e.preventDefault();

    const res = await fetch("http://localhost:8080/api/auth/change-password", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({ old_password: oldPassword, new_password: newPassword }),
    });

    const data = await res.json();
    if (res.ok) {
      setMessage("Пароль успешно изменён");
      setMessageColor("text-green-600");
      setOldPassword("");
      setNewPassword("");
    } else {
      setMessage(data.error || "Ошибка");
      setMessageColor("text-red-600");
    }
  };

  return (
    <form onSubmit={handleChange} className="max-w-md mx-auto mt-20 space-y-4 bg-white p-6 rounded shadow">
      <h2 className="text-xl font-bold">Смена пароля</h2>
      {message && <p className={`font-medium ${messageColor}`}>{message}</p>}
      <input
        type="password"
        placeholder="Старый пароль"
        className="w-full border p-2 rounded"
        value={oldPassword}
        onChange={(e) => setOldPassword(e.target.value)}
        required
      />
      <input
        type="password"
        placeholder="Новый пароль"
        className="w-full border p-2 rounded"
        value={newPassword}
        onChange={(e) => setNewPassword(e.target.value)}
        required
      />
      <button className="bg-blue-600 text-white px-4 py-2 rounded w-full">Изменить</button>
      <button type="button" onClick={handleBack} className="mt-2 text-blue-600 hover:underline">
        ← Назад
      </button>
    </form>
  );
}
