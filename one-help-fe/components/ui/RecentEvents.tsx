import React, { useState, useEffect } from "react";
import { View, Text, FlatList } from "react-native";
import { Link } from "expo-router";
import EventCard from "@/components/ui/EventSmallCard";

export default function RecentEvents() {
  const [events, setEvents] = useState<any[]>([]);

  useEffect(() => {
    setTimeout(() => {
      setEvents([
        {
          id: "1",
          date: "15",
          month: "БЕР",
          title: "Донація крові",
          location: "Лікарня Охматдит",
          color: "red",
        },
        {
          id: "2",
          date: "20",
          month: "БЕР",
          title: "Art Auction",
          location: "City Gallery",
          color: "purple",
        },
      ]);
    }, 1000);
  }, []);

  return (
    <View className="bg-white px-4">
      <View className="flex-row justify-between items-center mb-2">
        <Text className="text-2xl text-black">Найближчі події</Text>
        <Link href="/events" className="pt-4 py-2">
          <Text className="text-primary text-lg">Див усі</Text>
        </Link>
      </View>

      <FlatList
        data={events}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => <EventCard event={item} />}
        horizontal
        showsHorizontalScrollIndicator={false}
      />
    </View>
  );
}
