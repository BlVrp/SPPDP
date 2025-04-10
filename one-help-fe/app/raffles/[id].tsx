import React, { useEffect, useState } from "react";
import {
  View,
  Text,
  Image,
  TouchableOpacity,
  ScrollView,
  ActivityIndicator,
  Alert,
} from "react-native";
import { useLocalSearchParams, useRouter } from "expo-router";
import AsyncStorage from "@react-native-async-storage/async-storage";
import { Platform } from "react-native";
import * as FileSystem from "expo-file-system";
import * as Sharing from "expo-sharing";
import { useAuth } from "@/context/AuthContext";

export default function RaffleDetail() {
  const { id } = useLocalSearchParams();
  const router = useRouter();
  const { user } = useAuth();
  const [raffle, setRaffle] = useState<any>(null);
  const [fundraise, setFundraise] = useState<any>(null);
  const [participants, setParticipants] = useState<any[]>([]);
  const [showParticipants, setShowParticipants] = useState(false);
  const [currentGiftIndex, setCurrentGiftIndex] = useState(0);
  const [loading, setLoading] = useState(true);

  const fetchRaffle = async () => {
    try {
      const token = await AsyncStorage.getItem("token");
      if (!token) {
        Alert.alert("Помилка", "Користувач не авторизований");
        return;
      }

      const res = await fetch(`http://localhost:8080/api/v0/raffles/${id}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) throw new Error("Не вдалося отримати розіграш");
      const result = await res.json();
      setRaffle(result);

      const fundraiseRes = await fetch(
        `http://localhost:8080/api/v0/fundraises/${result.fundraiseId}`,
        { headers: { Authorization: `Bearer ${token}` } }
      );
      const fundraiseData = await fundraiseRes.json();
      setFundraise(fundraiseData);

      const partRes = await fetch(
        `http://localhost:8080/api/v0/users/raffle-participants/${id}`,
        { headers: { Authorization: `Bearer ${token}` } }
      );
      const participantsData = await partRes.json();
      setParticipants(participantsData);
    } catch (err: any) {
      Alert.alert("Помилка", err.message);
    } finally {
      setLoading(false);
    }
  };

  const downloadCSV = async () => {
    const csv = ["Ім'я,Прізвище,Місто"];
    participants.forEach((p) => {
      csv.push(`${p.firstName},${p.lastName},${p.city}`);
    });
    const csvString = csv.join("\n");

    if (Platform.OS === "web") {
      // WEB: Створюємо Blob і завантажуємо
      try {
        const blob = new Blob([csvString], { type: "text/csv" });
        const url = URL.createObjectURL(blob);
        const link = document.createElement("a");
        link.href = url;
        link.setAttribute("download", "participants.csv");
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
      } catch (error) {
        console.error("Помилка експорту на Web:", error);
        Alert.alert("Помилка", "Не вдалося завантажити CSV файл.");
      }
    } else {
      // MOBILE (iOS/Android): Через FileSystem + Sharing
      try {
        const fileUri = FileSystem.documentDirectory + "participants.csv";
        await FileSystem.writeAsStringAsync(fileUri, csvString, {
          encoding: FileSystem.EncodingType.UTF8,
        });
        await Sharing.shareAsync(fileUri);
      } catch (error) {
        console.error("Помилка експорту на мобільному:", error);
        Alert.alert("Помилка", "Не вдалося поділитися CSV файлом.");
      }
    }
  };

  useEffect(() => {
    if (id) fetchRaffle();
  }, [id]);

  const nextGift = () => {
    setCurrentGiftIndex((prev) =>
      prev + 1 < raffle.gifts.length ? prev + 1 : 0
    );
  };
  const prevGift = () => {
    setCurrentGiftIndex((prev) =>
      prev - 1 >= 0 ? prev - 1 : raffle.gifts.length - 1
    );
  };

  if (loading) {
    return (
      <View className="flex-1 justify-center items-center bg-white">
        <ActivityIndicator size="large" color="#2563EB" />
      </View>
    );
  }

  return (
    <ScrollView className="flex-1 bg-white p-4">
      <View className="bg-accent p-5 rounded-3xl shadow-sm">
        {raffle.gifts.length > 1 && (
          <View className="flex-row justify-center items-center space-x-8 mb-2">
            <TouchableOpacity onPress={prevGift}>
              <Text className="text-3xl">⬅</Text>
            </TouchableOpacity>
            <Text className="text-sm text-gray-500">
              {currentGiftIndex + 1} / {raffle.gifts.length}
            </Text>
            <TouchableOpacity onPress={nextGift}>
              <Text className="text-3xl">➡</Text>
            </TouchableOpacity>
          </View>
        )}

        <Image
          source={{ uri: raffle.gifts[currentGiftIndex].imageUrl }}
          className="w-full h-64 rounded-xl mb-2"
          resizeMode="cover"
        />

        <Text className="text-lg font-medium text-gray-700 text-center mb-2">
          🎁 {raffle.gifts[currentGiftIndex].title}
        </Text>

        <TouchableOpacity
          className="bg-primary rounded-full py-3 mt-4 items-center"
          onPress={() => router.push(`/fundraises/${raffle.fundraiseId}`)}
        >
          <Text className="text-white text-base font-semibold">
            Донат від {raffle.minimumDonation} ₴ 💰
          </Text>
        </TouchableOpacity>

        <Text className="text-2xl font-bold text-black text-center mt-6">
          {raffle.title}
        </Text>

        <Text className="text-gray-600 text-base text-center mt-3">
          {raffle.description}
        </Text>

        {/* Toggle participants
        <TouchableOpacity
          className="mt-4 bg-gray-200 py-2 px-4 rounded-full self-center"
          onPress={() => setShowParticipants((prev) => !prev)}
        >
          <Text className="text-center text-black font-medium">
            {showParticipants ? "Сховати учасників" : "Показати учасників"}
          </Text>
        </TouchableOpacity> */}

        {showParticipants &&
          participants.length > 0 &&
          user?.id === fundraise?.organizerId && (
            <View className="mt-4">
              <TouchableOpacity
                className="mt-4 bg-gray-200 py-2 px-4 rounded-full self-center"
                onPress={() => setShowParticipants((prev) => !prev)}
              >
                <Text className="text-center text-black font-medium">
                  {showParticipants
                    ? "Сховати учасників"
                    : "Показати учасників"}
                </Text>
              </TouchableOpacity>

              {participants.map((p, i) => (
                <View
                  key={i}
                  className="bg-white p-3 mb-2 rounded-xl shadow-sm"
                >
                  <Text className="text-black font-semibold">
                    {p.firstName} {p.lastName}
                  </Text>
                  <Text className="text-gray-500 text-sm">{p.city}</Text>
                </View>
              ))}

              {user?.id === fundraise?.organizerId && (
                <TouchableOpacity
                  className="mt-6 bg-blue-600 py-2 px-4 rounded-full self-center"
                  onPress={downloadCSV}
                >
                  <Text className="text-white font-medium">
                    ⬇️ Завантажити CSV
                  </Text>
                </TouchableOpacity>
              )}
            </View>
          )}
      </View>
    </ScrollView>
  );
}
