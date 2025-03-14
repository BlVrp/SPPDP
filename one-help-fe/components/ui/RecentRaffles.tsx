import React, { useState, useEffect } from "react";
import { View, Text, FlatList } from "react-native";
import { Link } from "expo-router";
import SmallRaffleCard from "@/components/ui/RaffleSmallCard";

export default function ActiveRafflesList() {
  const [raffles, setRaffles] = useState<any[]>([]);

  useEffect(() => {
    setTimeout(() => {
      setRaffles([
        {
          id: "1",
          image:
            "https://soundmagcdn.fra1.cdn.digitaloceanspaces.com/news/816/byZvgJY1WNcsfWAs.webp",
          title: "Ліжка для війська",
          description:
            "2 принти картини «Свята земля» з автографами KALUSH ORCHESTRA",
          minDonation: 100,
        },
        {
          id: "2",
          image: "https://s.dou.ua/img/announces/dou_5.png",
          title: "Допомога фронту",
          description: "Унікальна футболка з автографами",
          minDonation: 150,
        },
      ]);
    }, 1000);
  }, []);

  return (
    <View className="bg-white p-4">
      <View className="flex-row justify-between items-center mb-4">
        <Text className="text-2xl text-black">Активні розіграші</Text>
        <Link href="/raffles" className="pt-4 py-2">
          <Text className="text-primary text-lg">Див усі</Text>
        </Link>
      </View>

      <FlatList
        data={raffles}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => <SmallRaffleCard raffle={item} />}
        horizontal
        showsHorizontalScrollIndicator={false}
      />
    </View>
  );
}
