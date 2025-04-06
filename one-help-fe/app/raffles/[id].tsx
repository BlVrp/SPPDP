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
import { ProgressBar } from "react-native-paper";

export default function DetailedFundraiseCard() {
  const { id } = useLocalSearchParams();
  const router = useRouter();
  const [fundraiser, setFundraiser] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  const fetchFundraiser = async () => {
    try {
      const token = await AsyncStorage.getItem("token");
      if (!token) {
        Alert.alert("Помилка", "Користувач не авторизований");
        return;
      }

      const response = await fetch(`http://localhost:8080/api/v0/fundraises/${id}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || "Не вдалося отримати дані збору");
      }

      const data = await response.json();
      setFundraiser(data);
    } catch (err: any) {
      Alert.alert("Помилка", err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (id) fetchFundraiser();
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

  return (
    <ScrollView className="flex-1 bg-white p-4">
      <View className="bg-accent p-4 rounded-2xl">
        <Image
          source={{
            uri:
              fundraiser.image ||
              "https://via.placeholder.com/600x400.png?text=Fundraise",
          }}
          className="w-full h-96 rounded-lg mb-4"
          resizeMode="cover"
        />

        <Text className="text-xl font-bold text-black text-center mb-2">
          {fundraiser.title}
        </Text>

        <Text className="text-gray-msg mt-2 text-base leading-5 mb-4">
          {fundraiser.description}
        </Text>

        <View className="mt-4">
          <Text className="text-grey-msg text-sm text-center font-medium">
            {fundraiser.filledAmount.toLocaleString()} /{" "}
            {fundraiser.targetAmount.toLocaleString()} грн
          </Text>
          <ProgressBar
            progress={
              fundraiser.filledAmount / fundraiser.targetAmount || 0
            }
            color="#2563EB"
            className="h-3 rounded-md mt-1"
          />
        </View>

        <TouchableOpacity className="bg-primary rounded-lg p-3 mt-5 items-center">
          <Text className="text-white text-lg font-semibold">Донат 🍩</Text>
        </TouchableOpacity>
      </View>
    </ScrollView>
  );
}
