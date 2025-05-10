import React, { useEffect, useState } from "react";
import { View, Text, FlatList, ActivityIndicator, Alert } from "react-native";
import { Link } from "expo-router";
import AsyncStorage from "@react-native-async-storage/async-storage";
import EventCard from "@/components/ui/EventSmallCard";

export default function RecentEvents() {
  const [events, setEvents] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  const formatDateParts = (isoDate: string) => {
    const date = new Date(isoDate);
    const day = String(date.getDate()).padStart(2, "0");
    const month = date
      .toLocaleString("uk-UA", { month: "short" })
      .toUpperCase();
    return { date: day, month };
  };

  useEffect(() => {
    const fetchEvents = async () => {
      try {
        const token = await AsyncStorage.getItem("token");
        if (!token) {
          Alert.alert("Помилка", "Користувач не авторизований");
          return;
        }

        const res = await fetch(
          "http://localhost:8080/api/v0/events?limit=2&page=1",
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );

        if (!res.ok) {
          const err = await res.json();
          throw new Error(err.message || "Не вдалося завантажити події");
        }

        const result = await res.json();
        const formatted = result.data.map((event: any) => {
          const { date, month } = formatDateParts(event.startDate);
          return {
            id: event.id,
            date,
            month,
            title: event.title,
            location:
              event.format === "ONLINE"
                ? "Онлайн"
                : event.address || "Без адреси",
            color: "blue", 
          };
        });

        setEvents(formatted);
      } catch (err: any) {
        Alert.alert("Помилка", err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchEvents();
  }, []);

  return (
    <View className="bg-white px-4">
      <View className="flex-row justify-between items-center mb-2">
        <Text className="text-2xl text-black">Найближчі події</Text>
        <Link href="/events" className="pt-4 py-2">
          <Text className="text-primary text-lg">Див усі</Text>
        </Link>
      </View>

      {loading ? (
        <ActivityIndicator size="small" color="#2563EB" />
      ) : (
        <FlatList
          data={events}
          keyExtractor={(item) => item.id}
          renderItem={({ item }) => <EventCard event={item} />}
          horizontal
          showsHorizontalScrollIndicator={false}
        />
      )}
    </View>
  );
}
