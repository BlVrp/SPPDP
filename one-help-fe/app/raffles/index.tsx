import React, { useState, useEffect } from "react";
import {
  View,
  Text,
  FlatList,
  ActivityIndicator,
  Alert,
  TouchableOpacity,
} from "react-native";
import { Link } from "expo-router";
import RaffleCard from "@/components/ui/RaffleCard";
import AsyncStorage from "@react-native-async-storage/async-storage";

export default function RafflesList() {
  const [raffles, setRaffles] = useState<any[]>([]);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);

  const LIMIT = 2;

  const fetchRaffles = async () => {
    try {
      setLoading(true);
  
      const token = await AsyncStorage.getItem("token");
      if (!token) {
        Alert.alert("Помилка", "Користувач не авторизований");
        return;
      }
  
      const currentRes = await fetch(
        `http://localhost:8080/api/v0/raffles?limit=${LIMIT}&page=${page}`,
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      );
  
      if (!currentRes.ok) {
        const error = await currentRes.json();
        throw new Error(error.message || "Помилка під час отримання розіграшів");
      }
  
      const currentData = await currentRes.json();
      setRaffles(currentData.data || []);
  
      const nextRes = await fetch(
        `http://localhost:8080/api/v0/raffles?limit=${LIMIT}&page=${page + 1}`,
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      );
  
      const nextData = await nextRes.json();
      setHasMore((nextData.data?.length || 0) > 0);
    } catch (err: any) {
      Alert.alert("Помилка", err.message);
    } finally {
      setLoading(false);
    }
  };
  

  useEffect(() => {
    fetchRaffles();
  }, [page]);

  return (
    <View className="flex-1 bg-white p-4">
      <View className="flex-row justify-between items-center mb-4">
        <Text className="text-2xl text-black">Активні Розіграші</Text>
        <Link href="/raffles/create" className="pt-4 py-2">
          <Text className="text-primary text-lg">+ Створити розіграш</Text>
        </Link>
      </View>

      {loading ? (
        <ActivityIndicator size="large" color="#2563EB" />
      ) : raffles.length > 0 ? (
        <>
          <FlatList
            data={raffles}
            keyExtractor={(item) => item.id}
            renderItem={({ item }) => (
              <RaffleCard
                raffle={{
                  raffle_id: item.id,
                  title: item.title,
                  description: item.description,
                  minimum_donation: item.minimumDonation,
                  start_date: item.startDate,
                  end_date: item.endDate,
                  image: item.gifts?.[0]?.imageUrl || "",
                }}
              />
            )}
            showsVerticalScrollIndicator={false}
          />

          {/* Pagination */}
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
          Розіграшів не знайдено.
        </Text>
      )}
    </View>
  );
}
