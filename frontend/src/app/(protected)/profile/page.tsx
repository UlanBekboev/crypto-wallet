"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { fetchWithRefresh } from "@/lib/fetchWithRefresh";
import Link from "next/link";
import { usePublicClient, useAccount, useBalance } from "wagmi";
import { formatEther } from "viem";
import { ConnectButton } from "@rainbow-me/rainbowkit";

type User = {
  id: number;
  email: string;
  name?: string;
};

export default function ProfilePage() {
  const [user, setUser] = useState<User | null>(null);
  const [ethBalance, setEthBalance] = useState<string | null>(null);
  const [isClient, setIsClient] = useState(false);

  const router = useRouter();
  const publicClient = usePublicClient();
  const { address, isConnected } = useAccount();
  const { data: balance, isLoading } = useBalance({ address });

  // Устанавливаем, что компонент смонтирован на клиенте
  useEffect(() => {
    setIsClient(true);
  }, []);

  // Получение профиля
  useEffect(() => {
    if (!isConnected) return; // Без подключения не запрашиваем профиль

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
  }, [isConnected]);

  // Получение баланса
  useEffect(() => {
    const fetchBalance = async () => {
      if (!address || !publicClient) return;
      try {
        const balance = await publicClient.getBalance({ address });
        setEthBalance(formatEther(balance));
      } catch (err) {
        console.error("Ошибка получения баланса:", err);
      }
    };
    fetchBalance();
  }, [address, publicClient]);

  // Пока не клиент — ничего не рендерим
  if (!isClient) return null;

  // Если не подключён кошелёк
  if (!isConnected) {
    return (
      <div className="flex flex-col items-center justify-center text-center mt-10">
        <p className="text-gray-600 mt-2">Пожалуйста, подключите кошелек</p>
        <ConnectButton />
      </div>
    );
  }

  // Если подключён, но профиль ещё грузится
  if (!user) return <div className="text-center mt-10">Загрузка профиля...</div>;

  return (
    <div className="max-w-xl mx-auto mt-20 bg-white p-6 rounded shadow space-y-4">
      <h2 className="text-2xl font-bold">Ваш профиль</h2>
      <p><strong>Адрес:</strong> {address}</p>
      <p>
        <strong>Баланс:</strong>{' '}
        {isLoading ? 'Загрузка...' : `${balance?.formatted} ${balance?.symbol}`}
      </p>
      <p><strong>ID:</strong> {user.id}</p>
      <p><strong>Email:</strong> {user.email}</p>
      {user.name && <p><strong>Имя:</strong> {user.name}</p>}

      <div className="pt-4 border-t">
        <h3 className="text-lg font-semibold mb-2">Связь с MetaMask</h3>
        <ConnectButton accountStatus="address" chainStatus="icon" showBalance />
        {address && ethBalance && (
          <p className="mt-2 text-sm text-gray-700">
            💰 Баланс ETH: <strong>{ethBalance}</strong>
          </p>
        )}
      </div>

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
