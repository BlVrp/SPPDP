import React, { useEffect, useState } from "react";
import { SafeAreaView, ScrollView } from "react-native";
import Welcome from "@/components/ui/Welcome";
import RecentFunsraises from "@/components/ui/RecentFunsraises";
import RecentEvents from "@/components/ui/RecentEvents";
import RecentRaffles from "@/components/ui/RecentRaffles";
import AsyncStorage from "@react-native-async-storage/async-storage";

export default function HomeScreen() {
  const [userName, setUserName] = useState("Користувач");

  useEffect(() => {
    const fetchUserInfo = async () => {
      try {
        const userInfoStr = await AsyncStorage.getItem("user");
        if (userInfoStr) {
          const userInfo = JSON.parse(userInfoStr);
          setUserName(userInfo?.firstName || "Користувач");
        }
      } catch (error) {
        console.log("❌ Error fetching user info", error);
      }
    };

    fetchUserInfo();
  }, []);

  return (
    <ScrollView
      bounces={false}
      overScrollMode="never"
      showsVerticalScrollIndicator={false}
    >
      <SafeAreaView className="flex-1 bg-white">
        <Welcome
          user={{
            name: userName,
            donated: Math.floor(Math.random() * 1500 + 50), // e.g., 50-550₴
            events: Math.floor(Math.random() * 5 + 1),     // e.g., 1-5 events
            raffles: Math.floor(Math.random() * 3),        // e.g., 0-2 raffles
          }}
        />
        <RecentFunsraises />
        <RecentEvents />
        <RecentRaffles />
      </SafeAreaView>
    </ScrollView>
  );
}
