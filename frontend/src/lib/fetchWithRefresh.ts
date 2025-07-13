export async function fetchWithRefresh(
  input: RequestInfo,
  init?: RequestInit
): Promise<Response> {
  const res = await fetch(input, {
    ...init,
    credentials: "include",
  });

  if (res.status !== 401) return res;

  // Попробуем обновить токен
  const refreshRes = await fetch("http://localhost:8080/api/auth/refresh", {
    method: "POST",
    credentials: "include",
  });

  if (!refreshRes.ok) {
    throw new Error("Не удалось обновить токен");
  }

  // Повторяем исходный запрос
  return fetch(input, {
    ...init,
    credentials: "include",
  });
}
