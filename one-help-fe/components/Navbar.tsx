import React from 'react';
import { View, Text, TouchableOpacity } from 'react-native';

const Navbar = () => {
  return (
    <View className="flex flex-row justify-between items-center p-4 bg-red-400 shadow-md">
      <Text className="text-xl font-bold">Logo</Text>
      <TouchableOpacity>
        <Text className="text-blue-500">Menu</Text>
      </TouchableOpacity>
    </View>
  );
};

export default Navbar;
