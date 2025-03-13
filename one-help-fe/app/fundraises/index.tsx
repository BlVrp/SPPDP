import React, { useState, useEffect } from "react";
import {
  View,
  Text,
  FlatList,
  TouchableOpacity,
  ActivityIndicator,
} from "react-native";
import { useNavigation } from "expo-router";
import FundraiserCard from "@/components/ui/FundraiserCard";
import { Link } from "expo-router";

export default function FundraisersList() {
  const navigation = useNavigation();

  // Тимчасові збори (замість API)
  const [fundraisers, setFundraisers] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Симулюємо затримку запиту до API
    setTimeout(() => {
      setFundraisers([
        {
          id: "1",
          title: "На 2 ecoflow delta max",
          description: "для батальйону радіоелектронної розвідки",
          image:
            "https://www.peoplesproject.com/wp-content/uploads/2023/09/11113.jpg",
          goal: 71000,
          raised: 55000,
        },
        {
          id: "2",
          title: "Збір на авто",
          description: "25 бригада, Покровського напрямку",
          image:
            "https://ptv.ua/mediafiles/gallery/d5d26a74-41b8-4985-bccb-492ffa84f7b5.jfif",
          goal: 260000,
          raised: 100000,
        },
      ]);
      setLoading(false);
    }, 1500);
  }, []);

  return (
    <View className="flex-1 bg-white p-4">
      {/* Заголовок */}
      <View className="flex-row justify-between items-center mb-4">
        <Text className="text-2xl text-black">Активні збори</Text>
        <Link href="/fundraises/create" className="pt-4 py-2">
          <Text className="text-primary text-lg">+ Створити збір</Text>
        </Link>
      </View>

      {/* Список зборів */}
      {loading ? (
        <ActivityIndicator size="large" color="#2563EB" />
      ) : (
        <FlatList
          data={fundraisers}
          keyExtractor={(item) => item.id}
          renderItem={({ item }) => <FundraiserCard fundraiser={item} />}
          showsVerticalScrollIndicator={false}
        />
      )}
    </View>
  );
}
