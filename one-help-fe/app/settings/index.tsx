import React, { useState } from 'react';
import { View, Text, TextInput, TouchableOpacity, Alert, ScrollView } from 'react-native';
import { useRouter } from 'expo-router';
import { useAuth } from "@/context/AuthContext";

interface User {
  firstName: string;
  lastName: string;
  phoneNumber: string;
  email: string;
  city: string;
  post: string;
  postDepartment: string;
}

const initialUser: User = {
  firstName: 'Олеся',
  lastName: 'Іваненко',
  phoneNumber: '+380501234567',
  email: 'olesia@example.com',
  city: 'Київ',
  post: 'NP',
  postDepartment: 'Відділення №12',
};

export default function SettingsPage() {
  const [user, setUser] = useState<User>(initialUser);
  const [isEditingBasic, setIsEditingBasic] = useState(false);
  const [isEditingAddress, setIsEditingAddress] = useState(false);
  const { isLoggedIn, logout } = useAuth();
  const router = useRouter();

  const handleEditToggleBasic = () => {
    setIsEditingBasic(!isEditingBasic);
  };

  const handleEditToggleAddress = () => {
    setIsEditingAddress(!isEditingAddress);
  };

  const handleSave = () => {
    if (!user.firstName || !user.lastName || !user.phoneNumber || !user.email || !user.city || !user.post || !user.postDepartment) {
      Alert.alert('Помилка', 'Усі поля обов\'язкові');
      return;
    }
    setIsEditingBasic(false);
    setIsEditingAddress(false);
    //Alert.alert('Успіх', 'Дані успішно збережено');
  };

  const handleChange = (field: keyof User, value: string) => {
    setUser((prev) => ({ ...prev, [field]: value }));
  };

  const handleLoginLogout = async () => {
    if (isLoggedIn) {
      await logout(); // clear token or user data
      router.replace("/auth/login"); // redirect
    } else {
      router.push("/auth/login"); // go to login page
    }
  };

  return (
    <ScrollView
  className="flex-1 bg-white"
  contentContainerStyle={{ padding: 24, paddingBottom: 30 }}
>
      <Text className="text-2xl font-bold text-gray-msg mb-4">Налаштування профілю</Text>

      {['firstName', 'lastName', 'phoneNumber', 'email'].map((field) => (
        <View key={field} className="mb-4">
          <Text className={`text-gray-msg font-medium ${isEditingBasic ? 'text-blue-600' : ''}`}>{
            field === 'firstName' ? 'Ім\'я' : field === 'lastName' ? 'Прізвище' : field === 'phoneNumber' ? 'Телефон' : 'Email'
          }</Text>
          <TextInput
            className={`border ${isEditingBasic ? 'border-blue-600' : 'border-gray-300'} rounded-lg p-3 mt-1 ${isEditingBasic ? 'text-black' : 'text-gray-500 bg-gray-100'}`}
            placeholder={field}
            value={user[field]}
            editable={isEditingBasic}
            onChangeText={(value) => handleChange(field as keyof User, value)}
          />
        </View>
      ))}

      <TouchableOpacity
        onPress={isEditingBasic ? handleSave : handleEditToggleBasic}
        className="bg-blue-600 rounded-xl py-1 mt-2 items-center w-full self-center"
      >
        <Text className="text-white text-base font-semibold">{isEditingBasic ? '💾 Зберегти' : '✏️ Редагувати'}</Text>
      </TouchableOpacity>

      <View className="mt-6 p-4 bg-gray-100 rounded-lg">
        <Text className="text-xl font-semibold text-gray-msg mb-2">Адреса доставки:</Text>
        {['city', 'post', 'postDepartment'].map((field) => (
          <View key={field} className="mb-4">
            <Text className={`text-gray-msg font-medium ${isEditingAddress ? 'text-blue-600' : ''}`}>{
              field === 'city' ? 'Місто' : field === 'post' ? 'Пошта' : 'Відділення'
            }</Text>
            <TextInput
              className={`border ${isEditingAddress ? 'border-blue-600' : 'border-gray-300'} rounded-lg p-3 mt-1 ${isEditingAddress ? 'text-black' : 'text-gray-500 bg-gray-100'}`}
              placeholder={field}
              value={user[field]}
              editable={isEditingAddress}
              onChangeText={(value) => handleChange(field as keyof User, value)}
            />
          </View>
        ))}
      </View>

      <TouchableOpacity
        onPress={isEditingAddress ? handleSave : handleEditToggleAddress}
        className="bg-blue-600 rounded-xl py-1 mt-2 items-center w-full self-center"
      >
        <Text className="text-white text-base font-semibold">{isEditingAddress ? '💾 Зберегти' : '✏️ Редагувати'}</Text>
      </TouchableOpacity>

      <View className="mt-10 items-end">
        <TouchableOpacity
          onPress={handleLoginLogout}
          className="bg-red-500 px-6 py-2 rounded-xl"
        >
          <Text className="text-white text-base font-semibold">
            {isLoggedIn ? '🚪 Вийти' : '🔑 Увійти'}
          </Text>
        </TouchableOpacity>
      </View>
    </ScrollView>
  );
}