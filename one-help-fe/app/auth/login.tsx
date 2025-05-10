import React, { useState } from "react";
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  Alert,
  ScrollView,
} from "react-native";
import { useRouter } from "expo-router";
import AsyncStorage from "@react-native-async-storage/async-storage";
import { useAuth } from "../../context/AuthContext";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const router = useRouter();
  const { setToken, setUser } = useAuth();

  const handleLogin = async () => {
    if (!email || !password) {
      Alert.alert("Помилка", "Введіть Email та пароль");
      return;
    }

    try {
      const response = await fetch("http://localhost:8080/api/v0/auth/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          identifier: email,
          password,
        }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || "Помилка входу");
      }

      const data = await response.json();
      await AsyncStorage.setItem("token", data.token);
      await AsyncStorage.setItem("user", JSON.stringify(data.user));
      setToken(data.token);
      setUser(data.user);

      Alert.alert("Вхід успішний");
      router.push("/");
    } catch (error: any) {
      Alert.alert("Помилка", error.message);
    }
  };

  return (
    <ScrollView
      contentContainerStyle={{
        flexGrow: 1,
        justifyContent: "center",
        padding: 40,
      }}
      className="bg-white"
    >
      <View className="items-center mb-6">
        <Text className="text-4xl font-bold text-blue-600">💙 OneHelp</Text>
        <Text className="text-2xl font-semibold text-gray-800 mt-2">Вхід</Text>
      </View>

      <TextInput
        placeholder="Email"
        value={email}
        onChangeText={setEmail}
        className="border border-gray-300 rounded-lg p-3 mt-2 bg-gray-100 text-gray-700"
      />
      <TextInput
        placeholder="Пароль"
        value={password}
        onChangeText={setPassword}
        secureTextEntry
        className="border border-gray-300 rounded-lg p-3 mt-2 bg-gray-100 text-gray-700"
      />

      <TouchableOpacity
        onPress={handleLogin}
        className="bg-blue-600 py-1 rounded-xl mt-6 items-center"
      >
        <Text className="text-white text-lg font-semibold">Увійти 🔐</Text>
      </TouchableOpacity>

      <TouchableOpacity
        onPress={() => router.push("/auth/register")}
        className="mt-4 self-center"
      >
        <Text className="text-blue-600 text-base font-medium">
          Не маєте акаунту? Зареєструватись
        </Text>
      </TouchableOpacity>
    </ScrollView>
  );
}
