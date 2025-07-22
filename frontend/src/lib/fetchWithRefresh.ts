export async function fetchWithRefresh(
  input: RequestInfo,
  init?: RequestInit
): Promise<Response> {
  const response = await fetch(input, {
    ...init,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers || {}),
    },
  });

  if (response.status !== 401) {
    return response;
  }

  // Попробуем обновить access-токен
  const refreshResponse = await fetch("http://localhost:8080/api/auth/refresh", {
    method: "POST",
    credentials: "include",
  });

  if (!refreshResponse.ok) {
    throw new Error("Не удалось обновить токен");
  }

  // Повторяем исходный запрос
  return fetch(input, {
    ...init,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers || {}),
    },
  });
}
