import React, { useState, useEffect } from "react";
import { View, Text, FlatList, ActivityIndicator } from "react-native";
import EventCard from "@/components/ui/EventCard";
import { Link } from "expo-router";

export default function EventList() {
  const [events, setEvents] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setTimeout(() => {
      setEvents([
        {
          id: "1",
          title: "Благодійний марафон",
          description: "Захід для збору коштів на дрон",
          start_date: "2025-04-05",
          end_date: "2025-04-06",
          format: "offline",
          image:
            "https://fartlek.com.ua/wp-content/uploads/2022/08/balakliia-run-2022-1140x640.jpg",
          max_participants: 100,
          minimum_donation: 50,
          address: "Київ, вул. Хрещатик, 1",
          status: "active",
        },
        {
          id: "2",
          title: "Онлайн-лекція про тактичну медицину",
          description: "Дізнайтеся про першу допомогу в бойових умовах",
          image: "https://i.ytimg.com/vi/s-0V76nC1S4/maxresdefault.jpg",
          start_date: "2025-04-10",
          format: "online",
          max_participants: 500,
          minimum_donation: 0,
          address: "https://zoom.us/example",
          status: "active",
        },
      ]);
      setLoading(false);
    }, 1500);
  }, []);

  return (
    <View className="flex-1 bg-white p-4">
      <View className="flex-row justify-between items-center mb-4">
        <Text className="text-2xl text-black">Майбутні заходи</Text>
        <Link href="/events/create" className="pt-4 py-2">
          <Text className="text-primary text-lg">+ Створити подію</Text>
        </Link>
      </View>
      {loading ? (
        <ActivityIndicator size="large" color="#2563EB" />
      ) : (
        <FlatList
          data={events}
          keyExtractor={(item) => item.id}
          renderItem={({ item }) => <EventCard event={item} />}
          showsVerticalScrollIndicator={false}
        />
      )}
    </View>
  );
}
