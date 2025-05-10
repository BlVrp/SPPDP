import React, { useEffect, useState } from "react";
import {
  View,
  Text,
  Image,
  TouchableOpacity,
  ScrollView,
  ActivityIndicator,
  Alert,
  Linking,
} from "react-native";
import { useLocalSearchParams, useRouter } from "expo-router";
import AsyncStorage from "@react-native-async-storage/async-storage";
import defaultImage from "@/assets/images/field.jpeg";
import EventSmallCard from "@/components/ui/EventSmallCard";
import RaffleSmallCard from "@/components/ui/RaffleSmallCard";

export default function DetailedFundraiseCard() {
  const { id } = useLocalSearchParams();
  const router = useRouter();
  const [fundraiser, setFundraiser] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [events, setEvents] = useState<any[]>([]);
  const [raffles, setRaffles] = useState<any[]>([]);

  const fetchFundraise = async () => {
    try {
      const token = await AsyncStorage.getItem("token");
      if (!token) {
        Alert.alert("Помилка", "Користувач не авторизований");
        return;
      }

      const [fundRes, eventsRes, rafflesRes] = await Promise.all([
        fetch(`http://localhost:8080/api/v0/fundraises/${id}`, {
          headers: { Authorization: `Bearer ${token}` },
        }),
        fetch(`http://localhost:8080/api/v0/events`, {
          headers: { Authorization: `Bearer ${token}` },
        }),
        fetch(`http://localhost:8080/api/v0/raffles`, {
          headers: { Authorization: `Bearer ${token}` },
        }),
      ]);

      const fundData = await fundRes.json();
      const allEvents = await eventsRes.json();
      const allRaffles = await rafflesRes.json();

      setEvents(allEvents?.data?.filter((e) => e.fundraiseId === id) || []);
      setRaffles(allRaffles?.data?.filter((r) => r.fundraiseId === id) || []);

      setFundraiser(fundData);
      setEvents(eventsData?.data || []);
      setRaffles(rafflesData?.data || []);
    } catch (err: any) {
      Alert.alert("Помилка", err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDonate = async () => {
    try {
      const token = await AsyncStorage.getItem("token");
      if (!token) {
        Alert.alert("Помилка", "Користувач не авторизований");
        return;
      }

      const res = await fetch(
        `http://localhost:8080/api/v0/fundraises/${id}/donate`,
        {
          method: "POST",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!res.ok) {
        const error = await res.json();
        throw new Error(error.message || "Не вдалося ініціювати донат");
      }

      const { paymentUrl } = await res.json();
      if (paymentUrl) {
        Linking.openURL(paymentUrl);
      } else {
        throw new Error("Посилання для оплати відсутнє");
      }
    } catch (err: any) {
      Alert.alert("Помилка", err.message);
    }
  };

  useEffect(() => {
    if (id) fetchFundraise();
  }, [id]);

  if (loading) {
    return (
      <View className="flex-1 justify-center items-center bg-white">
        <ActivityIndicator size="large" color="#2563EB" />
      </View>
    );
  }

  if (!fundraiser) {
    return (
      <View className="flex-1 justify-center items-center bg-white">
        <Text className="text-lg text-gray-msg">Збір не знайдено 😕</Text>
        <TouchableOpacity
          className="mt-4 px-4 py-2 bg-primary rounded-md"
          onPress={() => router.back()}
        >
          <Text className="text-white">Назад</Text>
        </TouchableOpacity>
      </View>
    );
  }

  const progress = Math.min(
    fundraiser.filledAmount / fundraiser.targetAmount || 0,
    1
  );

  return (
    <ScrollView className="flex-1 bg-[#f4f6fa] px-4 py-6">
      <View className="bg-white rounded-3xl shadow-md overflow-hidden">
        <Image
          source={
            fundraiser.imageUrl?.trim().length
              ? { uri: fundraiser.imageUrl.trim() }
              : defaultImage
          }
          className="w-full h-52"
          resizeMode="cover"
        />

        <View className="p-6">
          <Text className="text-2xl font-bold text-center text-black mb-4">
            {fundraiser.title}
          </Text>

          <Text className="text-gray-700 text-base text-center mb-4">
            {fundraiser.description}
          </Text>

          <Text className="text-sm text-center text-gray-500 mb-2">
            📅 Завершення:{" "}
            {new Date(fundraiser.endDate).toLocaleDateString("uk-UA")}
          </Text>

          <View className="mt-4">
            <View className="w-full bg-gray-200 h-4 rounded-full overflow-hidden">
              <View
                style={{ width: `${progress * 100}%` }}
                className="h-full bg-gradient-to-r from-blue-500 to-indigo-600 rounded-full"
              />
            </View>
            <Text className="text-center text-sm text-gray-600 mt-2">
              {fundraiser.filledAmount.toLocaleString()} /{" "}
              {fundraiser.targetAmount.toLocaleString()} грн
            </Text>
          </View>

          <TouchableOpacity
            onPress={handleDonate}
            className="bg-primary rounded-full mt-6 py-3 items-center shadow-md"
          >
            <Text className="text-white text-lg font-semibold">Донат 🍩</Text>
          </TouchableOpacity>
        </View>
      </View>

      {/* Connected Events */}
      {events.length > 0 && (
        <View className="mt-6">
          <Text className="text-xl font-semibold text-gray-800 mb-2 ml-1">
            Пов'язані події
          </Text>
          <ScrollView
            horizontal
            showsHorizontalScrollIndicator={false}
            className="mb-4"
          >
            {events.map((event, index) => {
              const date = new Date(event.startDate);
              return (
                <EventSmallCard
                  key={index}
                  event={{
                    id: event.id,
                    date: date.getDate().toString(),
                    month: date
                      .toLocaleString("uk-UA", { month: "short" })
                      .toUpperCase(),
                    title: event.title,
                    location: event.address || "Онлайн",
                    color: "blue",
                  }}
                />
              );
            })}
          </ScrollView>
        </View>
      )}

      {/* Connected Raffles */}
      {raffles.length > 0 && (
        <View className="mt-4">
          <Text className="text-xl font-semibold text-gray-800 mb-2 ml-1">
            Пов'язані розіграші
          </Text>
          <ScrollView
            horizontal
            showsHorizontalScrollIndicator={false}
            className="mb-4"
          >
            {raffles.map((raffle, index) => (
              <RaffleSmallCard key={index} raffle={raffle} />
            ))}
          </ScrollView>
        </View>
      )}
    </ScrollView>
  );
}
