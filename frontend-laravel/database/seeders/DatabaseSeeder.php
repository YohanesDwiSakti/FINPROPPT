<?php

namespace Database\Seeders;

use App\Models\User;
use Illuminate\Database\Console\Seeds\WithoutModelEvents;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\Hash;

class DatabaseSeeder extends Seeder
{
    use WithoutModelEvents;

    /**
     * Seed the application's database.
     */
    public function run(): void
    {
        User::updateOrCreate(
            ['email' => 'admin@tiki.test'],
            [
                'name' => 'Admin Hub Denpasar',
                'password' => Hash::make('admin123'),
                'role' => 'admin',
            ],
        );

        User::updateOrCreate(
            ['email' => 'customer@tiki.test'],
            [
                'name' => 'Customer Demo',
                'password' => Hash::make('customer123'),
                'role' => 'customer',
            ],
        );
    }
}
