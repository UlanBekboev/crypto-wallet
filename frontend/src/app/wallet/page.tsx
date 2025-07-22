"use client";

import { useEffect, useState } from "react";
import { fetchWithRefresh } from "@/lib/fetchWithRefresh";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { WalletIcon } from "lucide-react";

export default function WalletPage() {
  const [balance, setBalance] = useState<number | null>(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    const fetchWallet = async () => {
      try {
        const res = await fetchWithRefresh("http://localhost:8080/api/auth/wallet");
        const text = await res.text();
        console.log("Ответ от /wallet:", text);

        if (res.ok) {
          const data = JSON.parse(text);
          setBalance(data.wallet.balance);
        } else {
          console.error("Ошибка запроса: ", res.status);
          router.push("/login");
        }
      } catch (err) {
        console.error("Ошибка при загрузке кошелька:", err);
        router.push("/login");
      } finally {
        setLoading(false);
      }
    };

    fetchWallet();
  }, []);

  return (
    <div className="max-w-md mx-auto mt-16 p-6 bg-white rounded-2xl shadow-md text-center space-y-6">
      <h1 className="text-2xl font-bold text-gray-800">💼 Внутренний кошелек</h1>

      {loading ? (
        <p className="text-gray-500">Загрузка баланса...</p>
      ) : (
        <div className="bg-blue-50 border border-blue-300 p-6 rounded-xl shadow-inner">
          <div className="flex items-center justify-center gap-3 text-blue-700 text-xl font-semibold">
            <WalletIcon className="w-6 h-6" />
            <span>{balance} монет</span>
          </div>
        </div>
      )}

      <Link
        href="/profile"
        className="inline-block bg-blue-600 hover:bg-blue-700 transition text-white px-4 py-2 rounded-lg"
      >
        ← Назад в профиль
      </Link>
    </div>
  );
}
