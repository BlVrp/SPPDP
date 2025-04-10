import React, { useState, useEffect } from "react";
import {
  View,
  Text,
  FlatList,
  TouchableOpacity,
  ActivityIndicator,
  Alert,
} from "react-native";
import { Link } from "expo-router";
import AsyncStorage from "@react-native-async-storage/async-storage";
import FundraiserCard from "@/components/ui/FundraiserCard";

export default function FundraisersList() {
  const [fundraisers, setFundraisers] = useState<any[]>([]);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);

  const LIMIT = 2;

  const fetchFundraisers = async () => {
    try {
      setLoading(true);
      const token = await AsyncStorage.getItem("token");
  
      if (!token) {
        Alert.alert("Помилка", "Користувач не авторизований");
        return;
      }
  
      const response = await fetch(
        `http://localhost:8080/api/v0/fundraises?limit=${LIMIT}&page=${page}`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
  
      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || "Помилка отримання зборів");
      }
  
      const currentResult = await response.json();
      setFundraisers(currentResult.data || []);
  
      const nextResponse = await fetch(
        `http://localhost:8080/api/v0/fundraises?limit=${LIMIT}&page=${page + 1}`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
  
      const nextResult = await nextResponse.json();
      setHasMore((nextResult.data?.length || 0) > 0);
    } catch (err: any) {
      Alert.alert("Помилка", err.message);
    } finally {
      setLoading(false);
    }
  };
  

  useEffect(() => {
    fetchFundraisers();
  }, [page]);

  return (
    <View className="flex-1 bg-white p-4">
      <View className="flex-row justify-between items-center mb-4">
        <Text className="text-2xl text-black">Активні збори</Text>
        <Link href="/fundraises/create" className="pt-4 py-2">
          <Text className="text-primary text-lg">+ Створити збір</Text>
        </Link>
      </View>

      {loading ? (
        <ActivityIndicator size="large" color="#2563EB" />
      ) : fundraisers.length > 0 ? (
        <>
          <FlatList
            data={fundraisers}
            keyExtractor={(item) => item.id}
            renderItem={({ item }) => <FundraiserCard fundraiser={item} />}
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
          Зборів не знайдено.
        </Text>
      )}
    </View>
  );
}
