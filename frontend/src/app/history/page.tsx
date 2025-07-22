"use client";

import { useEffect, useState } from "react";
import { fetchWithRefresh } from "@/lib/fetchWithRefresh";
import Link from "next/link";

type Transaction = {
  id: number;
  from_user_id: number;
  to_user_id: number;
  amount: number;
  created_at: string;
};

export default function HistoryPage() {
  const [history, setHistory] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchHistory = async () => {
      try {
        const res = await fetchWithRefresh("http://localhost:8080/api/auth/history");
        if (res.ok) {
          const data = await res.json();
          setHistory(data.transactions);
        } else {
          console.error("Ошибка при загрузке истории");
        }
      } catch (err) {
        console.error("Ошибка:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchHistory();
  }, []);

  return (
    <div className="max-w-2xl mx-auto mt-10 p-6 bg-white rounded shadow">
      <h1 className="text-2xl font-bold mb-4 text-center">📜 История транзакций</h1>

      {loading ? (
        <p className="text-center">Загрузка...</p>
      ) : history.length === 0 ? (
        <p className="text-center">Нет транзакций</p>
      ) : (
        <ul className="space-y-3">
          {history.map((tx) => (
            <li key={tx.id} className="border p-3 rounded">
              <p>
                <strong>От:</strong> {tx.from_user_id} → <strong>Кому:</strong> {tx.to_user_id}
              </p>
              <p>
                💰 <strong>Сумма:</strong> {tx.amount}
              </p>
              <p className="text-sm text-gray-500">{new Date(tx.created_at).toLocaleString()}</p>
            </li>
          ))}
        </ul>
      )}

      <Link
        href="/profile"
        className="mt-6 inline-block bg-blue-600 text-white px-4 py-2 rounded"
      >
        ← Назад в профиль
      </Link>
    </div>
  );
}
