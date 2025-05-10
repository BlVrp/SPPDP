import React, { useState, useEffect } from "react";
import { View, Text, FlatList, ActivityIndicator, Alert } from "react-native";
import { Link } from "expo-router";
import FundraiseSmallCard from "@/components/ui/FundraiserSmallCard";
import AsyncStorage from "@react-native-async-storage/async-storage";

export default function RecentFunsraises() {
  const [fundraisers, setFundraisers] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchRecentFundraisers = async () => {
    try {
      setLoading(true);
      const token = await AsyncStorage.getItem("token");
      if (!token) {
        Alert.alert("Помилка", "Користувач не авторизований");
        return;
      }
  
      const response = await fetch(
        "http://localhost:8080/api/v0/fundraises",
        {
          method: "GET",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
  
      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || "Помилка отримання зборів");
      }
  
      const result = await response.json();
  
      const sorted = (result.data || [])
        .sort((a, b) => new Date(b.startDate).getTime() - new Date(a.startDate).getTime())
        .slice(0, 2);
  
      setFundraisers(sorted);
    } catch (error: any) {
      Alert.alert("Помилка", error.message);
    } finally {
      setLoading(false);
    }
  };
  

  useEffect(() => {
    fetchRecentFundraisers();
  }, []);

  return (
    <View className="bg-white px-4">
      <View className="flex-row justify-between items-center">
        <Text className="text-2xl text-black">Активні збори</Text>
        <Link href="/fundraises" className="pt-4 py-2">
          <Text className="text-primary text-lg">Див усі</Text>
        </Link>
      </View>

      {loading ? (
        <ActivityIndicator size="large" color="#2563EB" />
      ) : (
        <FlatList
          data={fundraisers}
          scrollEnabled={false}
          keyExtractor={(item) => item.id}
          renderItem={({ item }) => <FundraiseSmallCard fundraiser={item} />}
          showsVerticalScrollIndicator={false}
        />
      )}
    </View>
  );
}
