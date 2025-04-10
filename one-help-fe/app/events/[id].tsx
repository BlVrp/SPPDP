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
import { Linking } from "react-native";

export default function DetailedEventCard() {
  const { id } = useLocalSearchParams();
  const router = useRouter();

  const [event, setEvent] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  const formatDate = (date: string) =>
    new Date(date).toLocaleDateString("uk-UA");

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

  return (
    <ScrollView className="flex-1 bg-[#f4f6fa] px-4 py-6">
      <View className="bg-white rounded-3xl shadow-md overflow-hidden">
        <Image
          source={{
            uri:
              event.imageUrl?.trim().length > 0
                ? event.imageUrl
                : "https://me.usembassy.gov/wp-content/uploads/sites/250/Ukraine_grain_black_sea_Agriculture_3-1068x712-1-1068x684.jpg",
          }}
          className="w-full h-52"
          resizeMode="cover"
        />

        <View className="p-6">
          {/* Title */}
          <Text className="text-2xl font-bold text-center text-black mb-2">
            {event.title}
          </Text>

          {/* Description */}
          <Text className="text-gray-700 text-base text-center mb-4">
            {event.description}
          </Text>

          {/* Event Info */}
          <View className="bg-[#f9fafb] p-4 rounded-xl space-y-2">
            <Text className="text-gray-600 text-sm">
              📅 Дати: {formatDate(event.startDate)} –{" "}
              {event.endDate ? formatDate(event.endDate) : "—"}
            </Text>
            <Text className="text-gray-600 text-sm">
              📍 Місце:{" "}
              {event.format === "ONLINE"
                ? "Онлайн"
                : event.address || "Без адреси"}
            </Text>
            <Text className="text-gray-600 text-sm">
              👥 Учасників: {event.maxParticipants}
            </Text>
            <Text className="text-gray-600 text-sm">
              💰 Мін. внесок:{" "}
              {event.minimumDonation === 0
                ? "Безкоштовно"
                : `${event.minimumDonation} грн`}
            </Text>
          </View>

          {/* Button */}
          <TouchableOpacity
            className="bg-primary rounded-full py-3 mt-6 items-center shadow-md"
            onPress={() => {
              if (event.formUrl?.trim()) {
                Linking.openURL(event.formUrl);
              } else {
                console.log(
                  "Посилання відсутнє",
                  "Форма реєстрації наразі недоступна."
                );
              }
            }}
          >
            <Text className="text-white text-lg font-semibold">
              Взяти участь 🎟️
            </Text>
          </TouchableOpacity>
        </View>
      </View>
    </ScrollView>
  );
}
