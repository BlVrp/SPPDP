import React, { useEffect, useState } from 'react';
import { View, Text, TextInput, TouchableOpacity, Alert, ScrollView } from 'react-native';
import { useRouter } from 'expo-router';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { useAuth } from "@/context/AuthContext";

export default function SettingsPage() {
  const { isLoggedIn, logout, user, setUser } = useAuth();
  const [editedUser, setEditedUser] = useState(user);
  const [isEditingBasic, setIsEditingBasic] = useState(false);
  const [isEditingAddress, setIsEditingAddress] = useState(false);
  const router = useRouter();

  useEffect(() => {
    if (user) {
      setEditedUser(user);
    }
  }, [user]);

  const handleEditToggleBasic = () => {
    setIsEditingBasic(!isEditingBasic);
  };

  const handleEditToggleAddress = () => {
    setIsEditingAddress(!isEditingAddress);
  };

  const handleChange = (field: keyof typeof user, value: string) => {
    if (!editedUser) return;
    setEditedUser((prev) => prev ? { ...prev, [field]: value } : null);
  };

  const handleSave = async () => {
    if (!editedUser) return;
  
    const requiredFields = ['firstName', 'lastName', 'phoneNumber', 'email'];
    const allFieldsFilled = requiredFields.every(field => editedUser[field]);
  
    if (!allFieldsFilled) {
      Alert.alert('Помилка', 'Усі поля обов\'язкові');
      return;
    }
  
    try {
      const token = await AsyncStorage.getItem('token');
  
      const response = await fetch(`http://localhost:8080/api/v0/users/`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(editedUser),
      });
  
      if (!response.ok) {
        const err = await response.json();
        throw new Error(err.message || 'Помилка збереження даних');
      }
  
      const updated = await response.json();
      setUser(updated);
      await AsyncStorage.setItem('user', JSON.stringify(updated));
  
      setIsEditingBasic(false);
      setIsEditingAddress(false);
      Alert.alert('✅ Успіх', 'Інформацію оновлено');
    } catch (error: any) {
      console.error('❌ Error saving user:', error);
      Alert.alert('Помилка', error.message);
    }
  };
  
  

  const handleLoginLogout = async () => {
    if (isLoggedIn) {
      await logout();
      router.replace("/auth/login");
    } else {
      router.push("/auth/login");
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
          <Text className={`text-gray-msg font-medium ${isEditingBasic ? 'text-blue-600' : ''}`}>
            {field === 'firstName' ? 'Ім\'я' :
              field === 'lastName' ? 'Прізвище' :
              field === 'phoneNumber' ? 'Телефон' :
              'Email'}
          </Text>
          <TextInput
            className={`border ${isEditingBasic ? 'border-blue-600' : 'border-gray-300'} rounded-lg p-3 mt-1 ${isEditingBasic ? 'text-black' : 'text-gray-500 bg-gray-100'}`}
            placeholder={field}
            value={editedUser?.[field] || ''}
            editable={isEditingBasic}
            onChangeText={(value) => handleChange(field as keyof typeof user, value)}
          />
        </View>
      ))}

      <TouchableOpacity
        onPress={isEditingBasic ? handleSave : handleEditToggleBasic}
        className="bg-blue-600 rounded-xl py-1 mt-2 items-center w-full self-center"
      >
        <Text className="text-white text-base font-semibold">
          {isEditingBasic ? '💾 Зберегти' : '✏️ Редагувати'}
        </Text>
      </TouchableOpacity>

      <View className="mt-6 p-4 bg-gray-100 rounded-lg">
  <Text className="text-xl font-semibold text-gray-msg mb-2">Адреса доставки:</Text>

  {[
    { field: 'city', label: 'Місто' },
    { field: 'post', label: 'Пошта' },
    { field: 'postDepartment', label: 'Відділення' }
  ].map(({ field, label }) => (
    <View key={field} className="mb-4">
      <Text className={`text-gray-msg font-medium ${isEditingAddress ? 'text-blue-600' : ''}`}>
        {label}
      </Text>
      <TextInput
        className={`border ${isEditingAddress ? 'border-blue-600' : 'border-gray-300'} rounded-lg p-3 mt-1 ${isEditingAddress ? 'text-black' : 'text-gray-500 bg-gray-100'}`}
        placeholder={`Введіть ${label.toLowerCase()}`}
        value={editedUser?.[field as keyof typeof editedUser] ?? ''}
        editable={isEditingAddress}
        onChangeText={(value) => handleChange(field as keyof typeof editedUser, value)}
      />
    </View>
  ))}
</View>


      <TouchableOpacity
        onPress={isEditingAddress ? handleSave : handleEditToggleAddress}
        className="bg-blue-600 rounded-xl py-1 mt-2 items-center w-full self-center"
      >
        <Text className="text-white text-base font-semibold">
          {isEditingAddress ? '💾 Зберегти' : '✏️ Редагувати'}
        </Text>
      </TouchableOpacity>

      <View className="mt-5 items-end">
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
