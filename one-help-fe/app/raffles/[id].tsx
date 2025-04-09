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

export default function RaffleDetail() {
  const { id } = useLocalSearchParams();
  const router = useRouter();
  const [raffle, setRaffle] = useState<any>(null);
  const [currentGiftIndex, setCurrentGiftIndex] = useState(0);
  const [loading, setLoading] = useState(true);

  const fetchRaffle = async () => {
    try {
      const token = await AsyncStorage.getItem("token");
      if (!token) {
        Alert.alert("Помилка", "Користувач не авторизований");
        return;
      }

      const res = await fetch(`http://localhost:8080/api/v0/raffles/${id}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!res.ok) {
        const error = await res.json();
        throw new Error(error.message || "Не вдалося отримати розіграш");
      }

      const result = await res.json();
      setRaffle(result);
    } catch (err: any) {
      Alert.alert("Помилка", err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (id) fetchRaffle();
  }, [id]);

  const nextGift = () => {
    setCurrentGiftIndex((prev) =>
      prev + 1 < raffle.gifts.length ? prev + 1 : 0
    );
  };

  const prevGift = () => {
    setCurrentGiftIndex((prev) =>
      prev - 1 >= 0 ? prev - 1 : raffle.gifts.length - 1
    );
  };

  if (loading) {
    return (
      <View className="flex-1 justify-center items-center bg-white">
        <ActivityIndicator size="large" color="#2563EB" />
      </View>
    );
  }

  if (!raffle) {
    return (
      <View className="flex-1 justify-center items-center bg-white p-6">
        <Text className="text-lg text-gray-msg">Розіграш не знайдено 😕</Text>
      </View>
    );
  }

  return (
    <ScrollView className="flex-1 bg-white p-4">
      <View className="bg-accent p-5 rounded-3xl shadow-sm">
        {/* Gift Navigation */}
        {raffle.gifts.length > 1 && (
          <View className="flex-row justify-center items-center space-x-8 mb-2">
            <TouchableOpacity onPress={prevGift}>
              <Text className="text-3xl">⬅</Text>
            </TouchableOpacity>
            <Text className="text-sm text-gray-500">
              {currentGiftIndex + 1} / {raffle.gifts.length}
            </Text>
            <TouchableOpacity onPress={nextGift}>
              <Text className="text-3xl">➡</Text>
            </TouchableOpacity>
          </View>
        )}

        {/* Gift Image and Title */}
        {raffle.gifts.length > 0 && (
          <View className="items-center">
            <Image
              source={{ uri: raffle.gifts[currentGiftIndex].imageUrl }}
              className="w-full h-64 rounded-xl mb-2"
              resizeMode="cover"
            />

            <Text className="text-lg font-medium text-gray-700 text-center mb-2">
              🎁 {raffle.gifts[currentGiftIndex].title}
            </Text>
          </View>
        )}

        {/* Donate Button */}
        <TouchableOpacity
          className="bg-primary rounded-full py-3 mt-4 items-center"
          onPress={() => router.push(`/fundraises/${raffle.fundraiseId}`)}
        >
          <Text className="text-white text-base font-semibold">
            Донат від {raffle.minimumDonation} ₴ 💰
          </Text>
        </TouchableOpacity>

        {/* Title */}
        <Text className="text-2xl font-bold text-black text-center mt-6">
          {raffle.title}
        </Text>

        {/* Description */}
        <Text className="text-gray-600 text-base text-center mt-3">
          {raffle.description}
        </Text>
      </View>
    </ScrollView>
  );
}
