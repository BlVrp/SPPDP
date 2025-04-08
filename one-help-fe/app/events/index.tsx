import React, { useEffect, useState } from "react";
import {
  View,
  Text,
  FlatList,
  ActivityIndicator,
  TouchableOpacity,
  Alert,
} from "react-native";
import { Link } from "expo-router";
import AsyncStorage from "@react-native-async-storage/async-storage";
import EventCard from "@/components/ui/EventCard";
import defaultImage from "@/assets/images/field.jpeg";

export default function EventList() {
  const [events, setEvents] = useState<any[]>([]);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);

  const LIMIT = 5;

  const fetchEvents = async () => {
    try {
      setLoading(true);
      const token = await AsyncStorage.getItem("token");
      if (!token) {
        Alert.alert("Помилка", "Користувач не авторизований");
        return;
      }

      const res = await fetch(
        `http://localhost:8080/api/v0/events?limit=${LIMIT}&page=${page}`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!res.ok) {
        const error = await res.json();
        throw new Error(error.message || "Не вдалося отримати події");
      }

      const result = await res.json();
      const formatted = result.data.map((event: any) => ({
        ...event,
        start_date: formatDate(event.startDate),
        end_date: event.endDate ? formatDate(event.endDate) : null,
        format: event.format?.toLowerCase(),
        image: event.imageUrl ? event.imageUrl : "https://me.usembassy.gov/wp-content/uploads/sites/250/Ukraine_grain_black_sea_Agriculture_3-1068x712-1-1068x684.jpg",
        minimum_donation: event.minimumDonation,
      }));

      setEvents(formatted);

      const nextRes = await fetch(
        `http://localhost:8080/api/v0/events?limit=${LIMIT}&page=${page + 1}`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      const nextResult = await nextRes.json();
      setHasMore((nextResult.data?.length || 0) > 0);
    } catch (err: any) {
      Alert.alert("Помилка", err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchEvents();
  }, [page]);

  const formatDate = (date: string) => {
    return new Date(date).toISOString().split("T")[0];
  };

  return (
    <View className="flex-1 bg-white p-4">
      <View className="flex-row justify-between items-center mb-4">
        <Text className="text-2xl text-black">Майбутні заходи</Text>
        <Link href="/events/create" className="pt-4 py-2">
          <Text className="text-primary text-lg">+ Створити подію</Text>
        </Link>
      </View>

      {loading ? (
        <ActivityIndicator size="large" color="#2563EB" />
      ) : events.length > 0 ? (
        <>
          <FlatList
            data={events}
            keyExtractor={(item) => item.id}
            renderItem={({ item }) => <EventCard event={item} />}
            showsVerticalScrollIndicator={false}
          />

          <View className="flex-row justify-between items-center mt-4 px-4">
            <TouchableOpacity
              onPress={() => setPage((prev) => Math.max(prev - 1, 1))}
              disabled={page === 1}
              className={`px-4 py-2 rounded-lg ${
                page === 1 ? "bg-gray-300" : "bg-blue-600"
              }`}
            >
              <Text className="text-white font-semibold">⬅️ Назад</Text>
            </TouchableOpacity>

            <Text className="text-base font-medium">Сторінка {page}</Text>

            <TouchableOpacity
              onPress={() => setPage((prev) => prev + 1)}
              disabled={!hasMore}
              className={`px-4 py-2 rounded-lg ${
                !hasMore ? "bg-gray-300" : "bg-blue-600"
              }`}
            >
              <Text className="text-white font-semibold">Далі ➡️</Text>
            </TouchableOpacity>
          </View>
        </>
      ) : (
        <Text className="text-gray-500 text-center mt-10">
          Подій не знайдено.
        </Text>
      )}
    </View>
  );
}
