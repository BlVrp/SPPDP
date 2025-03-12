import { useState } from "react";
import { View, Text, TouchableOpacity } from "react-native";
import Icon from "react-native-vector-icons/Feather"; // Feather icons
import Sidebar from "./ui/SideBar";

export default function Navbar() {
  const [isSidebarOpen, setSidebarOpen] = useState(false);

  return (
    <View className="flex-row items-center justify-between bg-white p-4">
      {/* Burger Menu */}
      <TouchableOpacity onPress={() => setSidebarOpen(true)}>
        <Icon name="menu" size={28} color="black" />
      </TouchableOpacity>

      <View className="flex-row items-right space-x-2">
        <Text className="text-xl">💙</Text>
        <Text className="text-xl font-bold">OneHelp</Text>
      </View>

      {/* Sidebar Component */}
      <Sidebar isOpen={isSidebarOpen} onClose={() => setSidebarOpen(false)} />
    </View>
  );
}
