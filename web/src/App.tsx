import { useEffect, useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { api } from "@/lib/api";

export default function App() {
  const [message, setMessage] = useState("загрузка");

  useEffect(() => {
    let cancelled = false;
    api
      .GET("/api/health")
      .then(({ data, response }) => {
        if (cancelled) {
          return;
        }
        if (response.ok && data?.status) {
          setMessage(data.status);
          return;
        }
        setMessage("база недоступна");
      })
      .catch(() => {
        if (!cancelled) {
          setMessage("нет ответа");
        }
      });
    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <main className="flex min-h-svh items-center justify-center p-6">
      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle>Запись на звонок</CardTitle>
          <CardDescription>Каркас сервиса</CardDescription>
        </CardHeader>
        <CardContent>
          <p>Состояние API: {message}</p>
        </CardContent>
      </Card>
    </main>
  );
}
