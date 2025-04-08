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
import { useLocalSearchParams } from "expo-router";
import AsyncStorage from "@react-native-async-storage/async-storage";
import { useRouter } from "expo-router";

const router = useRouter();

export default function RaffleDetail() {
  const { id } = useLocalSearchParams();
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
      <View className="bg-accent p-4 rounded-2xl">

      {raffle.gifts.length > 0 && (
  <View className="relative items-center">
    {raffle.gifts.length > 1 && (
      <TouchableOpacity
        onPress={prevGift}
        className="absolute left-0 top-1/2 -translate-y-1/2 px-3 py-2"
      >
        <Text className="text-2xl">⬅</Text>
      </TouchableOpacity>
    )}

    <Image
      source={{ uri: raffle.gifts[currentGiftIndex].imageUrl }}
      className="w-3/4 h-72 rounded-lg"
      resizeMode="cover"
    />
    <Text className="text-lg text-grey-msg text-center mt-1">
      {raffle.gifts[currentGiftIndex].title}
    </Text>

    {raffle.gifts.length > 1 && (
      <TouchableOpacity
        onPress={nextGift}
        className="absolute right-0 top-1/2 -translate-y-1/2 px-3 py-2"
      >
        <Text className="text-2xl">➡</Text>
      </TouchableOpacity>
    )}
  </View>
)}


        <View className="mt-4 flex-row justify-center items-center">
        <TouchableOpacity
  className="bg-primary rounded-md px-4 w-2/3 py-2"
  onPress={() => router.push(`/fundraises/${raffle.fundraiseId}`)}
>
  <Text className="text-white text-lg font-semibold text-center">
    Донат від {raffle.minimumDonation} ₴ 💰
  </Text>
</TouchableOpacity>

          {/* <Text className="text-4xl ml-3">☑️</Text> */}
        </View>

        <Text className="text-2xl font-bold text-black text-center mt-4">
          {raffle.title}
        </Text>

        <View className="mb-3 mt-2">
          <Text className="text-grey-msg text-md text-center">
            {raffle.description}
          </Text>

          <Text className="text-grey-msg text-center mt-2">
            📆 Донати приймаються до{" "}
            {new Date(raffle.endDate).toLocaleDateString("uk-UA")}
          </Text>
        </View>
      </View>
    </ScrollView>
  );
}
