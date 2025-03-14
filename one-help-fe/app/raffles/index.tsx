import React from "react";
import { View, Text, FlatList } from "react-native";
import RaffleCard from "@/components/ui/RaffleCard";
import { Link } from "expo-router";

const TEMP_RAFFLES = [
  {
    raffle_id: "1",
    title: "Лотерея для підтримки дітей",
    description: "Підтримайте дітей, долучившись до благодійної лотереї!",
    minimum_donation: 50,
    start_date: "2025-04-01",
    end_date: "2025-04-15",
    image:
      "https://soundmagcdn.fra1.cdn.digitaloceanspaces.com/news/816/byZvgJY1WNcsfWAs.webp",
  },
  {
    raffle_id: "2",
    title: "Благодійний розіграш подарунків",
    description: "Беріть участь у розіграші та вигравайте круті призи!",
    minimum_donation: 100,
    start_date: "2025-05-01",
    end_date: "2025-05-10",
    image: "https://s.dou.ua/img/announces/dou_5.png",
  },
];

export default function RafflesList() {
  return (
    <View className="flex-1 bg-white p-4">
      <View className="flex-row justify-between items-center mb-4">
        <Text className="text-2xl text-black">Активні Розіграші</Text>
        <Link href="/raffles/create" className="pt-4 py-2">
          <Text className="text-primary text-lg">+ Створити розіграш</Text>
        </Link>
      </View>

      <FlatList
        data={TEMP_RAFFLES}
        keyExtractor={(item) => item.raffle_id}
        renderItem={({ item }) => <RaffleCard raffle={item} />}
      />
    </View>
  );
}
