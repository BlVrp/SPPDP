import React, { useEffect, useState } from "react";
import {
  View,
  Text,
  Image,
  TouchableOpacity,
  ScrollView,
  ActivityIndicator,
  Alert,
} from "react-native";
import { useLocalSearchParams, useRouter } from "expo-router";
import AsyncStorage from "@react-native-async-storage/async-storage";
import defaultImage from "@/assets/images/field.jpeg";

export default function DetailedEventCard() {
  const { id } = useLocalSearchParams();
  const router = useRouter();

  const [event, setEvent] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchEvent = async () => {
      try {
        const token = await AsyncStorage.getItem("token");
        if (!token) {
          Alert.alert("Помилка", "Користувач не авторизований");
          router.back();
          return;
        }

        const res = await fetch(`http://localhost:8080/api/v0/events/${id}`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (!res.ok) {
          const err = await res.json();
          throw new Error(err.message || "Не вдалося завантажити подію");
        }

        const data = await res.json();
        setEvent(data);
      } catch (err: any) {
        Alert.alert("Помилка", err.message);
        router.back();
      } finally {
        setLoading(false);
      }
    };

    fetchEvent();
  }, [id]);

  if (loading) {
    return (
      <View className="flex-1 justify-center items-center bg-white">
        <ActivityIndicator size="large" color="#2563EB" />
      </View>
    );
  }

  if (!event) {
    return (
      <View className="flex-1 justify-center items-center bg-white">
        <Text className="text-lg text-gray-msg">Подію не знайдено 😕</Text>
        <TouchableOpacity
          className="mt-4 px-4 py-2 bg-primary rounded-md"
          onPress={() => router.back()}
        >
          <Text className="text-white">Назад</Text>
        </TouchableOpacity>
      </View>
    );
  }

  const formatDate = (date: string) =>
    new Date(date).toISOString().split("T")[0];

  return (
    <ScrollView className="flex-1 bg-white p-4">
      <View className="bg-accent p-4 rounded-2xl">
        <Image
          source={{
            uri: event.imageUrl || "https://me.usembassy.gov/wp-content/uploads/sites/250/Ukraine_grain_black_sea_Agriculture_3-1068x712-1-1068x684.jpg",
          }}
          className="w-full h-96 rounded-lg mb-4"
          resizeMode="cover"
        />

        <Text className="text-xl font-bold text-black text-center mb-2">
          {event.title}
        </Text>

        <Text className="text-gray-msg mt-2 text-base leading-5">
          {event.description}
        </Text>

        <View className="mt-4 bg-white p-4 rounded-lg shadow-sm">
          <Text className="text-grey-msg">
            📅 {formatDate(event.startDate)} –{" "}
            {event.endDate ? formatDate(event.endDate) : "—"}
          </Text>
          <Text className="text-grey-msg">
            📍{" "}
            {event.format === "ONLINE"
              ? "Онлайн"
              : event.address || "Без адреси"}
          </Text>
          <Text className="text-grey-msg">
            👥 Макс. учасників: {event.maxParticipants}
          </Text>
          <Text className="text-grey-msg">
            💰 Мін. внесок:{" "}
            {event.minimumDonation === 0
              ? "Безкоштовно"
              : `${event.minimumDonation} грн`}
          </Text>
        </View>

        <TouchableOpacity className="bg-primary rounded-lg p-1 mt-5 items-center">
          <Text className="text-white text-lg font-semibold">
            Взяти участь 🎟️
          </Text>
        </TouchableOpacity>
      </View>
    </ScrollView>
  );
}
