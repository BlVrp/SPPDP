import React, { useState, useEffect } from "react";
import { View, Text, FlatList, ActivityIndicator, Alert } from "react-native";
import { Link } from "expo-router";
import AsyncStorage from "@react-native-async-storage/async-storage";
import SmallRaffleCard from "@/components/ui/RaffleSmallCard";

export default function ActiveRafflesList() {
  const [raffles, setRaffles] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchRaffles = async () => {
    try {
      const token = await AsyncStorage.getItem("token");
      if (!token) {
        Alert.alert("Помилка", "Користувач не авторизований");
        return;
      }

      const response = await fetch(
        "http://localhost:8080/api/v0/raffles",
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || "Не вдалося отримати розіграші");
      }

      const result = await response.json();

      const sorted = result.data
        .sort(
          (a: any, b: any) =>
            new Date(b.endDate).getTime() - new Date(a.endDate).getTime()
        )
        .slice(0, 2);

      setRaffles(sorted);
    } catch (err: any) {
      Alert.alert("Помилка", err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRaffles();
  }, []);

  return (
    <View className="bg-white p-4">
      <View className="flex-row justify-between items-center mb-4">
        <Text className="text-2xl text-black">Активні розіграші</Text>
        <Link href="/raffles" className="pt-4 py-2">
          <Text className="text-primary text-lg">Див усі</Text>
        </Link>
      </View>

      {loading ? (
        <ActivityIndicator size="large" color="#2563EB" />
      ) : (
        <FlatList
          data={raffles}
          keyExtractor={(item) => item.id}
          renderItem={({ item }) => <SmallRaffleCard raffle={item} />}
          horizontal
          showsHorizontalScrollIndicator={false}
        />
      )}
    </View>
  );
}
